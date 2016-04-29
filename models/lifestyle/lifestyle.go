package lifestyle

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
	_ "github.com/go-sql-driver/mysql"
)

type Lifestyle struct {
	ID          int       `json:"id,omitempty" xml:"id,omitempty"`
	DateAdded   time.Time `json:"dateAdded,omitempty" xml:"dateAdded,omitempty"`
	ParentID    int       `json:"parentID,omitempty" xml:"parentID,omitempty"`
	Name        string    `json:"name,omitempty" xml:"name,omitempty"`
	ShortDesc   string    `json:"shortDesc,omitempty" xml:"shortDesc,omitempty"`
	LongDesc    string    `json:"longDesc,omitempty" xml:"longDesc,omitempty"`
	Image       string    `json:"image,omitempty" xml:"image,omitempty"`
	IsLifestyle int       `json:"isLifestyle,omitempty" xml:"isLifestyle,omitempty"`
	Sort        int       `json:"sort,omitempty" xml:"sort,omitempty"`
	Contents    Contents  `json:"contents,omitempty" xml:"contents,omitempty"`
	Towables    Towables  `json:"towables,omitempty" xml:"towables,omitempty"`
}

type Lifestyles []Lifestyle

type Content struct {
	ID          int         `json:"id,omitempty" xml:"id,omitempty"`
	UserID      int         `json:"userID,omitempty" xml:"userID,omitempty"`
	Text        string      `json:"content,omitempty" xml:"content,omitempty"`
	ContentType ContentType `json:"contentType,omitempty" xml:"contentType,omitempty"`
	Deleted     bool        `json:"deleted,omitempty" xml:"deleted,omitempty"`
	PartID      int
}
type Contents []Content

type ContentType struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	HTML bool   `json:"html,omitempty" xml:"html,omitempty"`
}
type Towable struct {
	ID         int    `json:"id,omitempty" xml:"id,omitempty"`
	CatId      int    `json:"catId,omitempty" xml:"catId,omitempty"`
	Name       string `json:"name,omitempty" xml:"name,omitempty"`
	ShortDesc  string `json:"shortDesc,omitempty" xml:"shortDesc,omitempty"`
	Image      string `json:"image,omitempty" xml:"image,omitempty"`
	HitchClass string `json:"hitchClass,omitempty" xml:"hitchClass,omitempty"`
	TW         int    `json:"TW,omitempty" xml:"TW,omitempty"`
	GTW        int    `json:"GTW,omitempty" xml:"GTW,omitempty"`
	Message    string `json:"message,omitempty" xml:"message,omitempty"`
}
type Towables []Towable

var (
	getAllLifestyles = `select c.catID, c.catTitle, c.dateAdded, c.parentID,
							c.shortDesc, c.longDesc, c.image, c.isLifestyle,
							c.sort from Categories as c
							Join ApiKeyToBrand as akb on akb.brandID = c.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where c.isLifestyle = 1 && (ak.api_key = ? && (c.brandID = ? OR 0=?))
							order by c.sort`
	getLifestyle = `select
						c.catID, c.catTitle, c.dateAdded, c.parentID,
						c.shortDesc, c.longDesc, c.image, c.isLifestyle,
						c.sort
						from Categories as c
						Join ApiKeyToBrand as akb on akb.brandID = c.brandID
						Join ApiKey as ak on akb.keyID = ak.id
						where c.catID = ? && (ak.api_key = ? && (c.brandID = ? OR 0=?))
						limit 1`
	getLifestyleContent = `select ct.allowHTML, ct.type, c.text from Content as c
							join ContentBridge as cb on c.contentID = cb.contentID
							join ContentType as ct on c.cTypeID = ct.cTypeID
							where cb.catID = ?`
	getAllLifestyleContent = `select cb.catID, ct.allowHTML, ct.type, c.text from Content as c
							join ContentBridge as cb on c.contentID = cb.contentID
							join ContentType as ct on c.cTypeID = ct.cTypeID
							join Category as cat on cat.catID = cb.catID
							Join ApiKeyToBrand as akb on akb.brandID = cat.brandID
							Join ApiKey as ak on akb.keyID = ak.id
							where cb.catID > 0 && (ak.api_key = ? && (cat.brandID = ? OR 0=?))`
	getLifestyleTowables = `select
								t.trailerID, t.name, t.shortDesc, t.hitchClass, t.image, t.TW, t.GTW, t.message
								from Trailer as t
								join Lifestyle_Trailer as lt on t.trailerID = lt.trailerID
								where lt.catID = ?
								order by t.TW`

	getAllLifestyleTowables = `select
								t.trailerID, lt.catId, t.name, t.shortDesc, t.hitchClass, t.image, t.TW, t.GTW, t.message
								from Trailer as t
								join Lifestyle_Trailer as lt on t.trailerID = lt.trailerID
								join Category as cat on cat.catID = lt.catID
								Join ApiKeyToBrand as akb on akb.brandID = cat.brandID
								Join ApiKey as ak on akb.keyID = ak.id
								where (ak.api_key = ? && (cat.brandID = ? OR 0=?))
								order by t.TW`
	getContent = `SELECT c.contentID, c.text, c.cTypeID, c.userID, c.deleted, ct.type, ct.allowHTML FROM Content AS c LEFT JOIN ContentType AS ct ON ct.cTypeId = c.cTypeId WHERE c.contentID = ?`
	getTowable = `SELECT trailerID, image, name, TW, GTW, hitchClass, shortDesc, message FROM Trailer WHERE trailerID = ?`
)

func GetAll(ctx *middleware.APIContext) (ls Lifestyles, err error) {
	redis_key := "lifestyle:all:" + ctx.DataContext.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &ls)
		return ls, err
	}

	stmt, err := ctx.DB.Prepare(getAllLifestyles)
	if err != nil {
		return ls, err
	}
	defer stmt.Close()
	//get content and towables
	cs, err := getAllContent(ctx)
	contentMap := cs.ToMap()
	ts, err := getAllTowables(ctx)
	towMap := ts.ToMap()

	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	for res.Next() {
		var l Lifestyle
		err = res.Scan(&l.ID, &l.Name, &l.DateAdded, &l.ParentID, &l.ShortDesc, &l.LongDesc, &l.Image, &l.IsLifestyle, &l.Sort)
		if err != nil {
			return ls, err
		}
		//bind content and towables
		cChan := make(chan int)
		tChan := make(chan int)

		go func() {
			for _, val := range contentMap {
				if val.ID == l.ID {
					l.Contents = append(l.Contents, val)
				}
			}
			cChan <- 1
		}()

		go func() {
			for _, val := range towMap {
				if val.CatId == l.ID {
					l.Towables = append(l.Towables, val)
				}
			}
			tChan <- 1
		}()
		<-cChan
		<-tChan

		ls = append(ls, l)
	}
	defer res.Close()
	go redis.Setex(redis_key, ls, 86400)
	return ls, err
}

func (l *Lifestyle) Get(ctx *middleware.APIContext) (err error) {
	redis_key := "lifestyle:get:" + strconv.Itoa(l.ID) + ":" + ctx.DataContext.BrandString
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &l)
		return err
	}

	stmt, err := ctx.DB.Prepare(getLifestyle)
	if err != nil {
		return err
	}
	defer stmt.Close()

	err = stmt.QueryRow(l.ID, ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID).Scan(&l.ID, &l.Name, &l.DateAdded, &l.ParentID, &l.ShortDesc, &l.LongDesc, &l.Image, &l.IsLifestyle, &l.Sort)
	if err != nil {
		return err
	}
	err = l.contents(ctx)
	if err != nil {
		return err
	}
	err = l.towables(ctx)
	if err != nil {
		return err
	}
	go redis.Setex(redis_key, l, 86400)
	return nil
}

func getAllContent(ctx *middleware.APIContext) (cs Contents, err error) {

	stmt, err := ctx.DB.Prepare(getAllLifestyleContent)
	if err != nil {
		return cs, err
	}
	defer stmt.Close()

	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	for res.Next() {
		var c Content
		err = res.Scan(&c.ID, &c.ContentType.HTML, &c.ContentType.Name, &c.Text)
		if err != nil {
			return cs, err
		}
		cs = append(cs, c)
	}
	defer res.Close()
	return cs, err
}

func getAllTowables(ctx *middleware.APIContext) (ts Towables, err error) {

	stmt, err := ctx.DB.Prepare(getAllLifestyleTowables)
	if err != nil {
		return ts, err
	}
	defer stmt.Close()
	res, err := stmt.Query(ctx.DataContext.APIKey, ctx.DataContext.BrandID, ctx.DataContext.BrandID)
	for res.Next() {
		var t Towable
		err = res.Scan(&t.ID, &t.CatId, &t.Name, &t.ShortDesc, &t.HitchClass, &t.Image, &t.TW, &t.GTW, &t.Message)
		if err != nil {
			return ts, err
		}
		ts = append(ts, t)
	}
	defer res.Close()
	return ts, err
}

func (l *Lifestyle) contents(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(getLifestyleContent)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(l.ID)
	for res.Next() {
		var c Content
		err = res.Scan(&c.ContentType.HTML, &c.ContentType.Name, &c.Text)
		if err != nil {
			return err
		}
		l.Contents = append(l.Contents, c)
	}
	defer res.Close()
	return err
}

func (l *Lifestyle) towables(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(getLifestyleTowables)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Query(l.ID)
	for res.Next() {
		var t Towable
		err = res.Scan(&t.ID, &t.Name, &t.ShortDesc, &t.HitchClass, &t.Image, &t.TW, &t.GTW, &t.Message)
		if err != nil {
			return err
		}
		l.Towables = append(l.Towables, t)
	}
	defer res.Close()
	return err
}

func (c *Content) get(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(getContent)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(c.ID).Scan(&c.ID, &c.Text, &c.ContentType.ID, &c.UserID, &c.Deleted, &c.ContentType.Name, &c.ContentType.HTML)

	return err
}

func (t *Towable) get(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(getTowable)
	if err != nil {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(t.ID).Scan(&t.ID, &t.Image, &t.Name, &t.TW, &t.GTW, &t.HitchClass, &t.ShortDesc, &t.Message)

	return err
}
