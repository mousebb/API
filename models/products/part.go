package products

import (
	"github.com/curt-labs/GoAPI/helpers/apicontext"
	"github.com/curt-labs/GoAPI/helpers/database"
	"github.com/curt-labs/GoAPI/helpers/redis"
	"github.com/curt-labs/GoAPI/helpers/rest"
	"github.com/curt-labs/GoAPI/helpers/sortutil"
	"github.com/curt-labs/GoAPI/models/customer"
	"github.com/curt-labs/GoAPI/models/customer/content"
	"github.com/curt-labs/GoAPI/models/vehicle"
	"github.com/curt-labs/GoAPI/models/video"
	_ "github.com/go-sql-driver/mysql"
	"sync"

	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	GetPaginatedPartNumbers = `select distinct p.partID
                               from Part as p
                               join ApiKeyToBrand as akb on akb.brandID = p.brandID
							   join ApiKey as ak on ak.id = akb.keyID
                               where (p.status = 800 || p.status = 900)
                               && ak.api_key = ? && (p.brandID = ? or 0 = ?)
                               order by p.partID
                               limit ?,?`
	GetFeaturedParts = `select distinct p.partID
                        from Part as p
                        join ApiKeyToBrand as akb on akb.brandID = p.brandID
						join ApiKey as ak on ak.id = akb.keyID
                        where (p.status = 800 || p.status = 900) && p.featured = 1
                        && ak.api_key = ? && (p.brandID = ? or 0 = ?)
                        order by p.dateAdded desc
                        limit 0, ?`
	GetLatestParts = `select distinct p.partID
                      from Part as p
                      join ApiKeyToBrand as akb on akb.brandID = p.brandID
					  join ApiKey as ak on ak.id = akb.keyID
                      where (p.status = 800 || p.status = 900)
                      && ak.api_key = ? && (p.brandID = ? or 0 = ?)
                      order by p.dateAdded desc
                      limit 0,?`
	SubCategoryIDStmt = `select distinct cp.partID
                         from CatPart as cp
                         join Part as p on cp.partID = p.partID
                         where cp.catID IN(%s) and (p.status = 800 || p.status = 900)
                         order by cp.partID
                         limit %d, %d`
	basicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.oldPartNumber, p.partID, p.priceCode, pc.class, pc.image, p.brandID
                from Part as p
                left join Class as pc on p.classID = pc.classID
                where p.partID = ? && p.status in (800,900) limit 1`
	relatedPartStmt = `select distinct relatedID from RelatedPart
                where partID = ?
                order by relatedID`
	partContentStmt = `select ct.cTypeID, ct.type, ct.allowHTML, con.text
                from Content as con
                join ContentBridge as cb on con.contentID = cb.contentID
                join ContentType as ct on con.cTypeID = ct.cTypeID
                where cb.partID = ? && LOWER(ct.type) != 'appguide' && con.deleted = 0
                order by ct.type`
	partInstallSheetStmt = `select c.text from ContentBridge as cb
                    join Content as c on cb.contentID = c.contentID
                    join ContentType as ct on c.cTypeID = ct.cTypeID
                    where partID = ? && ct.type = 'InstallationSheet'
                    limit 1`

	getPartByOldpartNumber = `select partID, status, dateModified, dateAdded, shortDesc, priceCode, classID, featured, ACESPartTypeID, brandID from Part where oldPartNumber = ?`

	getAllPartsBasicsStmt = `select p.status, p.dateAdded, p.dateModified, p.shortDesc, p.oldPartNumber, p.partID, p.priceCode, pc.class, pc.image, p.brandID, pa.value as upc
								from Part as p
								left join Class as pc on p.classID = pc.classID
								left join PartAttribute as pa on pa.partID = p.partID
								join ApiKeyToBrand as akb on akb.brandID = p.brandID
								join ApiKey as ak on ak.id = akb.keyID
								where (p.status in (800,900) && pa.field = "upc") && (ak.api_key = ? && (p.brandID = 1 or 0 = ?))
								order by partID`
	//create
	createPart = `INSERT INTO Part (partID, status, dateAdded, shortDesc, oldPartNumber, priceCode, classID, featured, ACESPartTypeID,brandID)
                    VALUES(?,?,?,?,?,?,?,?,?,?)`
	createPartAttributeJoin = `INSERT INTO PartAttribute (partID, value, field, sort) VALUES (?,?,?,?)`
	createVehiclePartJoin   = `INSERT INTO VehiclePart (vehicleID, partID, drilling, exposed, installTime) VALUES (?,?,?,?,?)`
	createContentBridge     = `INSERT INTO ContentBridge (catID, partID, contentID) VALUES (?,?,?)`
	createRelatedPart       = `INSERT INTO RelatedPart (partID, relatedID) VALUES (?,?)`
	createPartCategoryJoin  = `INSERT INTO CatPart (catID, partID) VALUES (?,?)`

	//delete
	deletePart               = `DELETE FROM Part WHERE partID  = ?`
	deletePartAttributeJoins = `DELETE FROM PartAttribute WHERE partID = ?`
	deleteVehiclePartJoins   = `DELETE FROM VehiclePart WHERE partID = ?`
	deleteContentBridgeJoins = `DELETE FROM ContentBridge WHERE partID = ?`
	deleteRelatedParts       = `DELETE FROM RelatedPart WHERE partID = ?`
	deletePartCategoryJoins  = `DELETE FROM CatPart WHERE partID = ?`

	//update
	updatePart = `UPDATE Part SET status = ?, shortDesc = ?, priceCode = ?, classID = ?, featured = ?, ACESPartTypeID = ?, brandID = ? WHERE partID = ?`
)

type Part struct {
	ID                int               `json:"id" xml:"id,attr"`
	BrandID           int               `json:"brandId" xml:"brandId,attr"`
	Status            int               `json:"status" xml:"status,attr"`
	PriceCode         int               `json:"price_code" xml:"price_code,attr"`
	RelatedCount      int               `json:"related_count" xml:"related_count,attr"`
	AverageReview     float64           `json:"average_review" xml:"average_review,attr"`
	DateModified      time.Time         `json:"date_modified" xml:"date_modified,attr"`
	DateAdded         time.Time         `json:"date_added" xml:"date_added,attr"`
	ShortDesc         string            `json:"short_description" xml:"short_description,attr"`
	InstallSheet      *url.URL          `json:"install_sheet" xml:"install_sheet"`
	Attributes        []Attribute       `json:"attributes" xml:"attributes"`
	VehicleAttributes []string          `json:"vehicle_atttributes" xml:"vehicle_attributes"`
	Vehicles          []vehicle.Vehicle `json:"vehicles,omitempty" xml:"vehicles,omitempty"`
	Content           []Content         `json:"content" xml:"content"`
	Pricing           []Price           `json:"pricing" xml:"pricing"`
	Reviews           []Review          `json:"reviews" xml:"reviews"`
	Images            []Image           `json:"images" xml:"images"`
	Related           []int             `json:"related" xml:"related"`
	Categories        []Category        `json:"categories" xml:"categories"`
	Videos            []video.Video     `json:"videos" xml:"videos"`
	Packages          []Package         `json:"packages" xml:"packages"`
	Customer          CustomerPart      `json:"customer,omitempty" xml:"customer,omitempty"`
	Class             Class             `json:"class,omitempty" xml:"class,omitempty"`
	Featured          bool              `json:"featured,omitempty" xml:"featured,omitempty"`
	AcesPartTypeID    int               `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty"`
	Installations     Installations     `json:"installation,omitempty" xml:"installation,omitempty"`
	Inventory         PartInventory     `json:"inventory,omitempty" xml:"inventory,omitempty"`
	OldPartNumber     string            `json:"oldPartNumber,omitempty" xml:"oldPartNumber,omitempty"`
	UPC               string            `json:"upc,omitempty" xml:"upc,omitempty"`
}

type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
}

type PaginatedProductListing struct {
	Parts         []Part `json:"parts" xml:"parts"`
	TotalItems    int    `json:"total_items" xml:"total_items,attr"`
	ReturnedCount int    `json:"returned_count" xml:"returned_count,attr"`
	Page          int    `json:"page" xml:"page,attr"`
	PerPage       int    `json:"per_page" xml:"per_page,attr"`
	TotalPages    int    `json:"total_pages" xml:"total_pages,attr"`
}

type Class struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty"`
	Name  string `json:"name,omitempty" xml:"name,omitempty"`
	Image string `json:"image,omitempty" xml:"image,omitempty"`
}

type Installation struct { //aka VehiclePart Table
	ID          int             `json:"id,omitempty" xml:"id,omitempty"`
	Vehicle     vehicle.Vehicle `json:"vehicle,omitempty" xml:"vehicle,omitempty"`
	Part        Part            `json:"part,omitempty" xml:"part,omitempty"`
	Drilling    string          `json:"drilling,omitempty" xml:"v,omitempty"`
	Exposed     string          `json:"exposed,omitempty" xml:"exposed,omitempty"`
	InstallTime int             `json:"installTime,omitempty" xml:"installTime,omitempty"`
}

type Installations []Installation

func (p *Part) FromDatabase(dtx *apicontext.DataContext, omit ...string) error {
	var errs []string

	omitFromResponse := make(map[string]string)
	for _, o := range omit {
		omitArray := strings.Split(o, ",")
		for _, oa := range omitArray {
			omitFromResponse[oa] = oa
		}
	}

	var wg sync.WaitGroup
	var cats []Category
	var attrs []Attribute
	var prices []Price
	var revs []Review
	var avgRev float64
	var imgs []Image
	var vids []video.Video
	var related []int
	var pkgs []Package
	var cons []Content
	wg.Add(9)

	go func(tmp *Part) {
		if _, ok := omitFromResponse["attributes"]; !ok {
			attrErr := p.GetAttributes(dtx)
			if attrErr != nil {
				errs = append(errs, attrErr.Error())
			}
			attrs = p.Attributes
		}
		wg.Done()
	}(p)

	go func() {
		if _, ok := omitFromResponse["pricing"]; !ok {
			priceErr := p.GetPricing(dtx)
			if priceErr != nil {
				errs = append(errs, priceErr.Error())
			}
			prices = p.Pricing
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["reviews"]; !ok {
			reviewErr := p.GetActiveApprovedReviews(dtx)
			if reviewErr != nil {
				errs = append(errs, reviewErr.Error())
			}
			revs = p.Reviews
			avgRev = p.AverageReview
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["images"]; !ok {
			imgErr := p.GetImages(dtx)
			if imgErr != nil {
				errs = append(errs, imgErr.Error())
			}
			imgs = p.Images
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["videos"]; !ok {
			var vidErr error
			vids, vidErr = video.GetPartVideos(p.ID)
			if vidErr != nil {
				errs = append(errs, vidErr.Error())
			}
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["related"]; !ok {
			relErr := p.GetRelated(dtx)
			if relErr != nil {
				errs = append(errs, relErr.Error())
			}
			related = p.Related
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["packaging"]; !ok {
			pkgErr := p.GetPartPackaging(dtx)
			if pkgErr != nil {
				errs = append(errs, pkgErr.Error())
			}
			pkgs = p.Packages
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["category"]; !ok {
			p.Categories = make([]Category, 0)
			catErr := p.PartBreadcrumbs(dtx)
			if catErr != nil {
				errs = append(errs, catErr.Error())
			}
			cats = p.Categories
		}
		wg.Done()
	}()

	go func() {
		if _, ok := omitFromResponse["content"]; !ok {
			conErr := p.GetContent(dtx)
			if conErr != nil {
				errs = append(errs, conErr.Error())
			}
			cons = p.Content
		}
		wg.Done()
	}()
	var basicErr error
	if basicErr = p.Basics(dtx); basicErr != nil {
		if basicErr == sql.ErrNoRows {
			basicErr = errors.New("Part #" + strconv.Itoa(p.ID) + " does not exist.")
		}
		errs = append(errs, basicErr.Error())
	}

	wg.Wait()

	p.Categories = cats
	p.Videos = vids
	p.Attributes = attrs
	p.Pricing = prices
	p.Images = imgs
	p.Packages = pkgs
	p.Content = cons
	p.Related = related
	p.RelatedCount = len(related)
	p.Reviews = revs
	p.AverageReview = avgRev

	if basicErr != nil {
		return errors.New("Could not find part: " + basicErr.Error())
	}

	go func(tmp Part) {
		redis.Setex(fmt.Sprintf("part:%d:%s", tmp.ID, dtx.BrandString), tmp, redis.CacheTimeout)
	}(*p)

	return nil
}

func (p *Part) Get(dtx *apicontext.DataContext, omit ...string) error {
	var omitFromResponse string
	for i, o := range omit {
		if i != 0 {
			omitFromResponse += ","
		}
		omitFromResponse += o
	}
	var err error
	var custPart CustomerPart
	var pi PartInventory
	customerChan := make(chan int)

	go func(api_key string) {
		custPart = p.BindCustomer(dtx)
		pi, _ = p.GetInventory(api_key, "", dtx)

		customerChan <- 1
	}(dtx.APIKey)

	redis_key := fmt.Sprintf("part:%d:%s", p.ID, dtx.BrandString)

	part_bytes, err := redis.Get(redis_key)
	if len(part_bytes) > 0 && err == nil {
		json.Unmarshal(part_bytes, &p)
	}

	p.Status = 0
	if p.Status == 0 {
		if err := p.FromDatabase(dtx, omitFromResponse); err != nil {
			return err
		}
	}

	<-customerChan
	close(customerChan)

	p.Customer = custPart
	p.Inventory = pi

	return nil
}

func All(page, count int, dtx *apicontext.DataContext) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetPaginatedPartNumbers)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, page, count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}
		go func(id int) {
			p := Part{ID: id}
			p.Get(dtx)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}
	defer rows.Close()

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.AscByField(parts, "ID")

	return parts, nil
}

func Featured(count int, dtx *apicontext.DataContext) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetFeaturedParts)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}

		go func(id int) {
			p := Part{ID: id}
			p.Get(dtx)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}
	defer rows.Close()

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.DescByField(parts, "DateAdded")

	return parts, nil
}

func GetAllPartsBasics(dtx *apicontext.DataContext) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAllPartsBasicsStmt)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID)
	if err != nil {
		return parts, err
	}

	for rows.Next() {
		var p Part
		var short, price, class, oldNum, image []byte
		var upc *string
		err = rows.Scan(
			&p.Status,
			&p.DateAdded,
			&p.DateModified,
			&short,
			&oldNum,
			&p.ID,
			&price,
			&class,
			&image,
			&p.BrandID,
			&upc,
		)
		if err != nil {
			return parts, err
		}
		if short != nil {
			p.ShortDesc = string(short[:])
		}
		if upc != nil {
			p.UPC = *upc
		}
		if price != nil {
			p.PriceCode, err = strconv.Atoi(string(price[:]))
			if err != nil {
				return parts, err
			}
		}
		if class != nil {
			p.Class.Name = string(class[:])
		}
		if image != nil {
			p.Class.Image = string(image[:])
		}
		if oldNum != nil {
			p.OldPartNumber = string(oldNum[:])
		}
		parts = append(parts, p)
	}
	defer rows.Close()

	return parts, nil
}

func Latest(count int, dtx *apicontext.DataContext) ([]Part, error) {
	parts := make([]Part, 0)

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return parts, err
	}
	defer db.Close()

	stmt, err := db.Prepare(GetLatestParts)
	if err != nil {
		return parts, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(dtx.APIKey, dtx.BrandID, dtx.BrandID, count)
	if err != nil {
		return parts, err
	}

	iter := 0
	partChan := make(chan int)
	for rows.Next() {
		var partID int
		if err = rows.Scan(&partID); err != nil {
			return parts, err
		}

		go func(id int) {
			p := Part{ID: id}
			p.Get(dtx)
			parts = append(parts, p)
			partChan <- 1
		}(partID)
		iter++
	}
	defer rows.Close()

	for i := 0; i < iter; i++ {
		<-partChan
	}

	sortutil.DescByField(parts, "DateAdded")

	return parts, nil
}

func (p *Part) GetWithVehicle(vehicle *vehicle.Vehicle, api_key string, dtx *apicontext.DataContext) error {
	var errs []string

	superChan := make(chan int)
	noteChan := make(chan int)
	go func(key string) {
		p.Get(dtx)
		superChan <- 1
	}(api_key)
	go func() {
		notes, nErr := vehicle.GetNotes(p.ID)
		if nErr != nil && notes != nil {
			errs = append(errs, nErr.Error())
			p.VehicleAttributes = []string{}
		} else {
			p.VehicleAttributes = notes
		}
		noteChan <- 1
	}()

	<-superChan
	<-noteChan

	if len(errs) > 0 {
		return errors.New("Error: " + strings.Join(errs, ", "))
	}
	return nil
}

func (p *Part) GetById(id int, key string, dtx *apicontext.DataContext) {
	p.ID = id

	p.Get(dtx)
}

func (p *Part) Basics(dtx *apicontext.DataContext) error {

	redis_key := fmt.Sprintf("part:%d:basics:%s", p.ID, dtx.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p); err == nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	qry, err := db.Prepare(basicsStmt)
	if err != nil {
		return err
	}
	defer qry.Close()

	row := qry.QueryRow(p.ID)
	if row == nil {
		return errors.New("No Part Found for:" + string(p.ID))
	}

	var short, price, class, oldNum, image []byte
	err = row.Scan(
		&p.Status,
		&p.DateAdded,
		&p.DateModified,
		&short,
		&oldNum,
		&p.ID,
		&price,
		&class,
		&image,
		&p.BrandID,
	)
	if err != nil {
		return err
	}
	if short != nil {
		p.ShortDesc = string(short[:])
	}
	if price != nil {
		p.PriceCode, err = strconv.Atoi(string(price[:]))
		if err != nil {
			return err
		}
	}
	if class != nil {
		p.Class.Name = string(class[:])
	}
	if image != nil {
		p.Class.Image = string(image[:])
	}
	if oldNum != nil {
		p.OldPartNumber = string(oldNum[:])
	}
	if dtx.BrandString != "" {
		go func(tmp Part) {
			redis.Setex(redis_key, tmp, redis.CacheTimeout)
		}(*p)
	}

	return nil
}

func (p *Part) GetRelated(dtx *apicontext.DataContext) error {
	redis_key := fmt.Sprintf("part:%d:related:%s", p.ID, dtx.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Related); err == nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(relatedPartStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.ID)
	if err != nil {
		return err
	}

	var related []int
	var relatedID int
	for rows.Next() {
		err = rows.Scan(&relatedID)
		if err != nil {
			return err
		}
		related = append(related, relatedID)
	}
	defer rows.Close()

	p.Related = related
	p.RelatedCount = len(related)

	if dtx.BrandString != "" {
		go func(rel []int) {
			redis.Setex(redis_key, rel, redis.CacheTimeout)
		}(p.Related)
	}

	return nil
}

func (p *Part) GetContent(dtx *apicontext.DataContext) error {
	redis_key_content := fmt.Sprintf("part:%d:content:%s", p.ID, dtx.BrandString)
	redis_key_installSheet := fmt.Sprintf("part:%d:installSheet:%s", p.ID, dtx.BrandString)

	data_content, err_content := redis.Get(redis_key_content)
	data_installSheet, err_installSheet := redis.Get(redis_key_installSheet)
	if err_content == nil && err_installSheet == nil && len(data_content) > 0 && len(data_installSheet) > 0 {
		err_content = json.Unmarshal(data_content, &p.Content)
		err_installSheet = json.Unmarshal(data_installSheet, &p.InstallSheet)

		if err_content == nil && err_installSheet == nil {
			return nil
		}
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(partContentStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(p.ID)
	if err != nil {
		return err
	}

	var content []Content
	for rows.Next() {
		var con Content
		var conText *string
		err = rows.Scan(
			&con.ContentType.Id,
			&con.ContentType.Type,
			&con.ContentType.AllowHtml,
			&conText,
		)
		if err != nil {
			return err
		}
		if conText != nil {
			con.Text = *conText
		}

		if strings.Contains(strings.ToLower(con.ContentType.Type), "install") {
			//sheetUrl, _ := url.Parse(con.Value)
			p.InstallSheet, _ = url.Parse(con.Text)
			// p.InstallSheet, _ = url.Parse(api_helpers.API_DOMAIN + "/part/" + strconv.Itoa(p.ID) + ".pdf")
		} else {
			content = append(content, con)
		}
	}
	defer rows.Close()

	p.Content = content
	if dtx.BrandString != "" && len(p.Content) > 0 {
		go redis.Setex(redis_key_content, p.Content, redis.CacheTimeout)
	}
	if dtx.BrandString != "" && p.InstallSheet != nil {
		go redis.Setex(redis_key_installSheet, p.InstallSheet, redis.CacheTimeout)
	}
	return nil
}

func (p *Part) BindCustomer(dtx *apicontext.DataContext) CustomerPart {
	var price float64
	var ref int

	priceChan := make(chan int)
	refChan := make(chan int)
	contentChan := make(chan int)

	go func() {
		price, _ = customer.GetCustomerPrice(dtx, p.ID)
		priceChan <- 1
	}()

	go func() {
		ref, _ = customer.GetCustomerCartReference(dtx.APIKey, p.ID)
		refChan <- 1
	}()

	go func() {
		content, _ := custcontent.GetPartContent(p.ID, dtx.APIKey)
		for _, con := range content {

			strArr := strings.Split(con.ContentType.Type, ":")
			cType := con.ContentType.Type
			if len(strArr) > 1 {
				cType = strArr[1]
			}
			var c Content
			c.ContentType.Type = cType
			c.Text = con.Text
			p.Content = append(p.Content, c)
		}
		contentChan <- 1
	}()

	<-priceChan
	<-refChan
	<-contentChan

	return CustomerPart{
		Price:         price,
		CartReference: ref,
	}
}

func (p *Part) GetInstallSheet(r *http.Request, dtx *apicontext.DataContext) (data []byte, err error) {
	redis_key := fmt.Sprintf("part:%d:installSheet:%s", p.ID, dtx.BrandString)

	data, err = redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		return data, nil
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare(partInstallSheetStmt)
	if err != nil {
		return
	}
	defer stmt.Close()

	var text string
	err = stmt.QueryRow(p.ID).Scan(
		&text,
	)
	if err != nil {
		return
	}

	data, err = rest.GetPDF(text, r)
	if err != nil {
		return
	}

	if dtx.BrandString != "" {
		go func(dt []byte) {
			redis.Setex(redis_key, dt, redis.CacheTimeout)
		}(data)
	}
	return
}

// PartBreacrumbs
//
// Description: Builds out Category breadcrumb array for the current part object.
//
// Inherited: part Part
// Returns: error
func (p *Part) PartBreadcrumbs(dtx *apicontext.DataContext) error {
	if p.ID == 0 {
		return errors.New("Invalid Part Number")
	}

	//check redis!
	redis_key := fmt.Sprintf("part:%d:breadcrumbs:%s", p.ID, dtx.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err = json.Unmarshal(data, &p.Categories); err == nil {
			return nil
		}
	}

	// Oh alright, let's talk with our database
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	partCategoryStmt, err := db.Prepare(PartCategoryStmt)
	if err != nil {
		return err
	}
	defer partCategoryStmt.Close()

	lookupCategoriesStmt, err := db.Prepare(CategoriesByBrandStmt)
	if err != nil {
		return err
	}
	defer lookupCategoriesStmt.Close()

	// Execute SQL Query against current ID
	catRow := partCategoryStmt.QueryRow(p.ID)
	if catRow == nil {
		return errors.New("No part found for " + string(p.ID))
	}

	ch := make(chan Category)
	go PopulateCategory(catRow, ch, dtx)
	initCat := <-ch
	close(ch)

	// Build thee lookup
	catLookup := make(map[int]Category)
	rows, err := lookupCategoriesStmt.Query(dtx.BrandID)
	if err != nil {
		return err
	}
	defer rows.Close()

	multiChan := make(chan []Category, 0)
	go PopulateCategoryMulti(rows, multiChan)
	cats := <-multiChan
	close(multiChan)

	for _, cat := range cats {
		catLookup[cat.ID] = cat
	}

	// Okay, let's put it together!
	var categories []Category
	categories = append(categories, initCat)

	nextParentID := initCat.ParentID
	for {
		if nextParentID == 0 {
			break
		}
		if c, found := catLookup[nextParentID]; found {
			nextParentID = c.ParentID
			categories = append(categories, c)
			continue
		}
		nextParentID = 0
	}

	p.Categories = categories
	if dtx.BrandString != "" {
		go func(cats []Category) {
			redis.Setex(redis_key, cats, redis.CacheTimeout)
		}(p.Categories)
	}

	return nil
}

func (p *Part) GetPartCategories(dtx *apicontext.DataContext) (cats []Category, err error) {
	redis_key := fmt.Sprintf("part:%d:categories:%s", p.ID, dtx.BrandString)

	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &cats); err == nil {
			return cats, nil
		}
	}

	if p.ID == 0 {
		return
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return
	}
	defer db.Close()

	qry, err := db.Prepare(PartAllCategoryStmt)
	if err != nil {
		return
	}
	defer qry.Close()

	// Execute SQL Query against current ID
	catRows, err := qry.Query(p.ID)
	if err != nil || catRows == nil { // Error occurred while executing query
		return
	}

	ch := make(chan []Category, 0)
	go PopulateCategoryMulti(catRows, ch)
	cats = <-ch

	for _, cat := range cats {

		contentChan := make(chan int)
		subChan := make(chan int)
		customerChan := make(chan int)

		c := Category{
			ID: cat.ID,
		}

		go func() {
			content, contentErr := c.GetContent()
			if contentErr == nil {
				cat.Content = content
			}
			contentChan <- 1
		}()

		go func() {
			subs, subErr := c.GetSubCategories(dtx)
			if subErr == nil {
				cat.SubCategories = subs
			}
			subChan <- 1
		}()

		go func() {
			content, _ := custcontent.GetCategoryContent(cat.ID, dtx.APIKey)
			for _, con := range content {
				strArr := strings.Split(con.ContentType.Type, ":")
				cType := con.ContentType.Type
				if len(strArr) > 1 {
					cType = strArr[1]
				}
				var catCon Content
				catCon.ContentType.Type = cType
				catCon.Text = con.Text
				cat.Content = append(cat.Content, catCon)
			}
			customerChan <- 1
		}()

		<-contentChan
		<-subChan
		<-customerChan

		cats = append(cats, cat)
	}
	if dtx.BrandString != "" {
		go func(cts []Category) {
			redis.Setex(redis_key, cts, redis.CacheTimeout)
		}(cats)
	}
	return
}

func (p *Part) GetPartByOldPartNumber(key string) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(getPartByOldpartNumber)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var shortDesc *string
	var priceCode, acesPartTypeID *int
	err = stmt.QueryRow(p.OldPartNumber).Scan(
		&p.ID,
		&p.Status,
		&p.DateModified,
		&p.DateAdded,
		&shortDesc,
		&priceCode,
		&p.Class.ID,
		&p.Featured,
		&acesPartTypeID,
		&p.BrandID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Invalid Part Number")
			return
		}
		return err
	}

	if shortDesc != nil {
		p.ShortDesc = *shortDesc
	}
	if priceCode != nil {
		p.PriceCode = *priceCode
	}
	if acesPartTypeID != nil {
		p.AcesPartTypeID = *acesPartTypeID
	}

	return nil
}

func (p *Part) Create(dtx *apicontext.DataContext) (err error) {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	p.DateAdded = time.Now()
	_, err = stmt.Exec(
		p.ID,
		p.Status,
		p.DateAdded,
		p.ShortDesc,
		p.OldPartNumber,
		p.PriceCode,
		p.Class.ID,
		p.Featured,
		p.AcesPartTypeID,
		p.BrandID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()

	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	priceChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)

	go func() (err error) {
		for _, attribute := range p.Attributes {
			err = p.CreatePartAttributeJoin(attribute, dtx)
			if err != nil {
				pajChan <- 1
				return err
			}
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		for _, installation := range p.Installations {
			err = p.CreateInstallation(installation, dtx)
			if err != nil {
				diChan <- 1
				return err
			}
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		for _, content := range p.Content {
			err = p.CreateContentBridge(p.Categories, content, dtx)
			if err != nil {
				dcbChan <- 1
				return err
			}
		}
		dcbChan <- 1
		return err
	}()
	go func() (err error) {
		for _, price := range p.Pricing {
			price.PartId = p.ID
			err = price.Create(dtx)
			if err != nil {
				priceChan <- 1
				return err
			}
		}
		priceChan <- 1
		return err
	}()
	go func() (err error) {
		for _, review := range p.Reviews {
			review.PartID = p.ID
			err = review.Create(dtx)
			if err != nil {
				revChan <- 1
				return err
			}
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		for _, image := range p.Images {
			image.PartID = p.ID
			err = image.Create(dtx)
			if err != nil {
				imageChan <- 1
				return err
			}
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		for _, related := range p.Related {
			err = p.CreateRelatedPart(related, dtx)
			if err != nil {
				relatedChan <- 1
				return err
			}
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		for _, category := range p.Categories {
			err = p.CreatePartCategoryJoin(category, dtx)
			if err != nil {
				pcjChan <- 1
				return err
			}
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		for _, video := range p.Videos {
			// video.PartID = p.ID
			err = video.CreateJoinPart(p.ID)
			// err = video.CreatePartVideo(dtx)
			if err != nil {
				videoChan <- 1
				return err
			}
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		for _, pack := range p.Packages {
			pack.PartID = p.ID
			err = pack.Create(dtx)
			if err != nil {
				packChan <- 1
				return err
			}
		}
		packChan <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-priceChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan

	return err
}

func (p *Part) Delete(dtx *apicontext.DataContext) (err error) {
	if p.ID == 0 {
		return errors.New("Part ID is zero.")
	}
	go redis.Delete(fmt.Sprintf("part:%d:%s", p.ID, dtx.BrandString))

	var price Price
	price.PartId = p.ID
	err = price.DeleteByPart(dtx)
	if err != nil {
		return err
	}

	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)

	go func() (err error) {
		err = p.DeletePartAttributeJoins(dtx)
		if err != nil {
			pajChan <- 1
			return err
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteInstallations(dtx)
		if err != nil {
			diChan <- 1
			return err
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteContentBridges(dtx)
		if err != nil {
			dcbChan <- 1
			return err
		}
		dcbChan <- 1
		return err
	}()

	go func() (err error) {
		var review Review
		review.PartID = p.ID
		err = review.Delete(dtx)
		if err != nil {
			revChan <- 1
			return err
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		var image Image
		image.PartID = p.ID
		err = image.DeleteByPart(dtx)
		if err != nil {
			imageChan <- 1
			return err
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteRelatedParts(dtx)
		if err != nil {
			relatedChan <- 1
			return err
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeletePartCategoryJoins(dtx)
		if err != nil {
			pcjChan <- 1
			return err
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		var v PartVideo
		v.PartID = p.ID
		err = v.DeleteByPart(dtx)
		if err != nil {
			videoChan <- 1
			return err
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		var pack Package
		pack.PartID = p.ID
		err = pack.DeleteByPart(dtx)
		if err != nil {
			packChan <- 1
			return err
		}
		packChan <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Part) Update(dtx *apicontext.DataContext) (err error) {
	if p.ID == 0 {
		return errors.New("Part ID is zero.")
	}
	go redis.Delete(fmt.Sprintf("part:%d:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(updatePart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.Status, p.ShortDesc, p.PriceCode, p.Class.ID, p.Featured, p.AcesPartTypeID, p.BrandID, p.ID)
	if err != nil {
		return err
	}
	//Refresh joins
	pajChan := make(chan int)
	diChan := make(chan int)
	dcbChan := make(chan int)
	priceChan := make(chan int)
	revChan := make(chan int)
	imageChan := make(chan int)
	relatedChan := make(chan int)
	pcjChan := make(chan int)
	videoChan := make(chan int)
	packChan := make(chan int)
	pajChanC := make(chan int)
	diChanC := make(chan int)
	dcbChanC := make(chan int)
	priceChanC := make(chan int)
	revChanC := make(chan int)
	imageChanC := make(chan int)
	relatedChanC := make(chan int)
	pcjChanC := make(chan int)
	videoChanC := make(chan int)
	packChanC := make(chan int)

	go func() (err error) {
		err = p.DeletePartAttributeJoins(dtx)
		if err != nil {
			pajChan <- 1
			return err
		}
		pajChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteInstallations(dtx)
		if err != nil {
			diChan <- 1
			return err
		}
		diChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteContentBridges(dtx)
		if err != nil {
			dcbChan <- 1
			return err
		}
		dcbChan <- 1
		return err
	}()
	go func() (err error) {
		var price Price
		price.PartId = p.ID
		err = price.DeleteByPart(dtx)
		if err != nil {
			priceChan <- 1
			return err
		}
		priceChan <- 1
		return err
	}()
	go func() (err error) {
		var review Review
		review.PartID = p.ID
		err = review.Delete(dtx)
		if err != nil {
			revChan <- 1
			return err
		}
		revChan <- 1
		return err
	}()
	go func() (err error) {
		var image Image
		image.PartID = p.ID
		err = image.DeleteByPart(dtx)
		if err != nil {
			imageChan <- 1
			return err
		}
		imageChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeleteRelatedParts(dtx)
		if err != nil {
			relatedChan <- 1
			return err
		}
		relatedChan <- 1
		return err
	}()
	go func() (err error) {
		err = p.DeletePartCategoryJoins(dtx)
		if err != nil {
			pcjChan <- 1
			return err
		}
		pcjChan <- 1
		return err
	}()
	go func() (err error) {
		var v PartVideo
		v.PartID = p.ID
		err = v.DeleteByPart(dtx)
		if err != nil {
			videoChan <- 1
			return err
		}
		videoChan <- 1
		return err
	}()
	go func() (err error) {
		var pack Package
		pack.PartID = p.ID
		err = pack.DeleteByPart(dtx)
		if err != nil {
			packChan <- 1
			return err
		}
		packChan <- 1
		return err
	}()

	go func() (err error) {
		for _, attribute := range p.Attributes {
			err = p.CreatePartAttributeJoin(attribute, dtx)
			if err != nil {
				pajChanC <- 1
				return err
			}
		}
		pajChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, installation := range p.Installations {
			err = p.CreateInstallation(installation, dtx)
			if err != nil {
				diChanC <- 1
				return err
			}
		}
		diChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, content := range p.Content {
			err = p.CreateContentBridge(p.Categories, content, dtx)
			if err != nil {
				dcbChanC <- 1
				return err
			}
		}
		dcbChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, price := range p.Pricing {
			price.PartId = p.ID
			err = price.Create(dtx)
			if err != nil {
				priceChanC <- 1
				return err
			}
		}
		priceChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, review := range p.Reviews {
			review.PartID = p.ID
			err = review.Create(dtx)
			if err != nil {
				revChanC <- 1
				return err
			}
		}
		revChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, image := range p.Images {
			image.PartID = p.ID
			err = image.Create(dtx)
			if err != nil {
				imageChanC <- 1
				return err
			}
		}
		imageChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, related := range p.Related {
			err = p.CreateRelatedPart(related, dtx)
			if err != nil {
				relatedChanC <- 1
				return err
			}
		}
		relatedChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, category := range p.Categories {
			err = p.CreatePartCategoryJoin(category, dtx)
			if err != nil {
				pcjChanC <- 1
				return err
			}
		}
		pcjChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, video := range p.Videos {
			// video.PartID = p.ID
			err = video.CreateJoinPart(p.ID)
			// err = video.CreatePartVideo(dtx)
			if err != nil {
				videoChanC <- 1
				return err
			}
		}
		videoChanC <- 1
		return err
	}()
	go func() (err error) {
		for _, pack := range p.Packages {
			pack.PartID = p.ID
			err = pack.Create(dtx)
			if err != nil {
				packChanC <- 1
				return err
			}
		}
		packChanC <- 1
		return err
	}()

	<-pajChan
	<-diChan
	<-dcbChan
	<-priceChan
	<-revChan
	<-imageChan
	<-relatedChan
	<-pcjChan
	<-videoChan
	<-packChan
	<-pajChanC
	<-diChanC
	<-dcbChanC
	<-priceChanC
	<-revChanC
	<-imageChanC
	<-relatedChanC
	<-pcjChanC
	<-videoChanC
	<-packChanC

	return err
}

//Join Creators
func (p *Part) CreatePartAttributeJoin(a Attribute, dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:attributes:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartAttributeJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, a.Value, a.Key, a.Sort)
	if err != nil {
		return err
	}
	return nil
}

//Creates "VehiclePart" Join, which also contains installation fields
func (p *Part) CreateInstallation(i Installation, dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:installation:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createVehiclePartJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(i.Vehicle.ID, p.ID, i.Drilling, i.Exposed, i.InstallTime)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	i.ID = int(id)
	return nil
}

func (p *Part) CreateContentBridge(cats []Category, c Content, dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:content:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	stmt, err := tx.Prepare(createContentBridge)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, cat := range cats {
		_, err = stmt.Exec(cat.ID, p.ID, c.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	return err
}

func (p *Part) CreateRelatedPart(relatedID int, dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:related:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createRelatedPart)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID, relatedID)
	if err != nil {
		return err
	}
	return nil
}

func (p *Part) CreatePartCategoryJoin(c Category, dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:categories:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(createPartCategoryJoin)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.ID, p.ID)
	if err != nil {
		return err
	}
	return nil
}

//delete Joins
func (p *Part) DeletePartAttributeJoins(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:attributes:%s", p.ID, dtx.BrandString))

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartAttributeJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteInstallations(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:installation:%s", p.ID, dtx.BrandString))

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteVehiclePartJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteContentBridges(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:content:%s", p.ID, dtx.BrandString))

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteContentBridgeJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeleteRelatedParts(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:related:%s", p.ID, dtx.BrandString))

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deleteRelatedParts)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
func (p *Part) DeletePartCategoryJoins(dtx *apicontext.DataContext) (err error) {
	go redis.Delete(fmt.Sprintf("part:%d:category:%s", p.ID, dtx.BrandString))
	go redis.Delete(fmt.Sprintf("part:%d:categories:%s", p.ID, dtx.BrandString))
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(deletePartCategoryJoins)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.ID)
	if err != nil {
		return err
	}
	return nil
}
