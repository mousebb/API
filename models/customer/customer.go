package customer

import (
	"github.com/curt-labs/API/helpers/conversions"
	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/geography"

	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude" xml:"latitude"`
	Longitude float64 `json:"longitude" xml:"longitude"`
}

type Customer struct {
	Id                  int                 `json:"id,omitempty" xml:"id,omitempty"`
	Name                string              `json:"name,omitempty" xml:"name,omitempty"`
	Email               string              `json:"email,omitempty" xml:"email,omitempty"`
	Address             string              `json:"address,omitempty" xml:"address,omitempty"`
	Address2            string              `json:"address2,omitempty" xml:"address2,omitempty"`
	City                string              `json:"city,omitempty" xml:"city,omitempty"`
	State               geography.State     `json:"state,omitempty" xml:"state,omitempty"`
	PostalCode          string              `json:"postalCode,omitempty" xml:"postalCode,omitempty"`
	Phone               string              `json:"phone,omitempty" xml:"phone,omitempty"`
	Fax                 string              `json:"fax,omitempty" xml:"fax,omitempty"`
	ContactPerson       string              `json:"contactPerson,omitempty" xml:"contactPerson,omitempty"`
	Latitude            float64             `json:"coords,latitude,omitempty" xml:"latitude,omitempty"`
	Longitude           float64             `json:"coords,longitude,omitempty" xml:"longitude,omitempty"`
	Website             url.URL             `json:"website,omitempty" xml:"website,omitempty"`
	Parent              *Customer           `json:"parent,omitempty" xml:"parent,omitempty"`
	SearchUrl           url.URL             `json:"searchUrl,omitempty" xml:"searchUrl,omitempty"`
	Logo                url.URL             `json:"logo,omitempty" xml:"logo,omitempty"`
	DealerType          DealerType          `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	DealerTier          DealerTier          `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	Locations           []CustomerLocation  `json:"locations,omitempty" xml:"locations,omitempty"`
	Users               []CustomerUser      `json:"users,omitempty" xml:"users,omitempty"`
	CustomerId          int                 `json:"customerId,omitempty" xml:"customerId,omitempty"`
	IsDummy             bool                `json:"isDummy,omitempty" xml:"isDummy,omitempty"`
	ELocalUrl           url.URL             `json:"eLocalUrl,omitempty" xml:"eLocalUrl,omitempty"`
	MapixCode           MapixCode           `json:"mapixCode,omitempty" xml:"mapixCode,omitempty"`
	ApiKey              string              `json:"apiKey,omitempty" xml:"apiKey,omitempty"`
	ShowWebsite         bool                `json:"showWebsite,omitempty" xml:"showWebsite,omitempty"`
	SalesRepresentative SalesRepresentative `json:"salesRepresentative,omitempty" xml:"salesRepresentative,omitempty"`
	BrandIDs            []int               `json:"brandIds,omitempty" xml:"brandIds,omitempty"`
	Accounts            []Account           `json:"accounts,omitempty" xml:"accounts,omitempty"`
	ShippingInfo        ShippingInfo        `json:"shippingInfo,omitempty" xml:"shippingInfo,omitempty"`
}

type Customers []Customer

type Scanner interface {
	Scan(...interface{}) error
}

type Account struct {
	ID                 int          `json:"id,omitempty" xml:"id,omitempty"`
	AccountNumber      string       `json:"accountNumber,omitempty" xml:"accountNumber,omitempty"`
	Cust_id            int          `json:"cust_id,omitempty" xml:"cust_id,omitempty"`
	TypeID             int          `json:"-" xml:"-"`
	FreightLimit       float64      `json:"freightLimit,omitempty" xml:"freightLimit,omitempty"`
	DefaultWarehouseID int          `json:"defaultWarehouseId,omitempty" xml:"defaultWarehouseId,omitempty"`
	Type               AccountType  `json:"type,omitempty" xml:"type,omitempty"`
	ShippingInfo       ShippingInfo `json:"shipping_info,omitempty" xml"shipping_info,omitempty"`
}

type AccountType struct {
	ID        int     `json:"id,omitempty" xml:"id,omitempty"`
	Title     string  `json:"title,omitempty" xml:"title,omitempty"`
	ComnetURL url.URL `json:"comnetURL,omitempty" xml:"comnetURL,omitempty"`
}

type SalesRepresentative struct {
	ID   int    `json:"id,omitempty" xml:"id,omitempty"`
	Name string `json:"name,omitempty" xml:"name,omitempty"`
	Code string `json:"code,omitempty" xml:"code,omitempty"`
}

type DealerType struct {
	Id      int     `json:"id,omitempty" xml:"id,omitempty"`
	Type    string  `json:"type,omitempty" xml:"type,omitempty"`
	Label   string  `json:"label,omitempty" xml:"label,omitempty"`
	Online  bool    `json:"online,omitempty" xml:"online,omitempty"`
	Show    bool    `json:"show,omitempty" xml:"show,omitempty"`
	MapIcon MapIcon `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
}

type DealerTier struct {
	Id      int    `json:"id,omitempty" xml:"id,omitempty"`
	Tier    string `json:"tier,omitempty" xml:"tier,omitempty"`
	Sort    int    `json:"sort,omitempty" xml:"sort,omitempty"`
	BrandID int    `json:"brandId,omitempty" xml:"brandId,omitempty"`
}

type MapIcon struct {
	Id            int `json:"id,omitempty" xml:"id,omitempty"`
	TierId        int `json:"tierId,omitempty" xml:"tierId,omitempty"`
	DealerTypeId  int
	MapIcon       url.URL `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
	MapIconShadow url.URL `json:"mapIconShadow,omitempty" xml:"mapIconShadow,omitempty"`
}

type MapGraphics struct {
	DealerTier DealerTier `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	DealerType DealerType `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	MapIcon    MapIcon    `json:"mapIcon,omitempty" xml:"mapIcon,omitempty"`
}

type GeoLocation struct {
	Latitude  float64 `json:"latitude,omitempty" xml:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty" xml:"longitude,omitempty"`
}

type MapixCode struct {
	ID          int    `json:"id,omitempty" xml:"id,omitempty"`
	Code        string `json:"code,omitempty" xml:"code,omitempty"`
	Description string `json:"description,omitempty" xml:"description,omitempty"`
}

type DealerLocation struct {
	CustomerLocation
	Distance            float64             `json:"distance,omitempty" xml:"distance,omitempty"`
	Website             url.URL             `json:"website,omitempty" xml:"website,omitempty"`
	Parent              *Customer           `json:"parent,omitempty" xml:"parent,omitempty"`
	SearchUrl           url.URL             `json:"searchUrl,omitempty" xml:"searchUrl,omitempty"`
	Logo                url.URL             `json:"logo,omitempty" xml:"logo,omitempty"`
	DealerType          DealerType          `json:"dealerType,omitempty" xml:"dealerType,omitempty"`
	DealerTier          DealerTier          `json:"dealerTier,omitempty" xml:"dealerTier,omitempty"`
	SalesRepresentative SalesRepresentative `json:"salesRepresentative,omitempty" xml:"salesRepresentative,omitempty"`
	MapixCode           MapixCode           `json:"mapixCode,omitempty" xml:"mapixCode,omitempty"`
}
type DealerLocations []DealerLocation

type StateRegion struct {
	Id           int          `json:"id,omitempty" xml:"id,omitempty"`
	Name         string       `json:"name,omitempty" xml:"name,omitempty"`
	Abbreviation string       `json:"abbreviation,omitempty" xml:"abbreviation,omitempty"`
	Count        int          `json:"count,omitempty" xml:"count,omitempty"`
	Polygons     []MapPolygon `json:"polygon,omitempty" xml:"polygon,omitempty"`
}

type MapPolygon struct {
	Id          int           `json:"id,omitempty" xml:"id,omitempty"`
	Coordinates []GeoLocation `json:"coordinates,omitempty" xml:"coordinates,omitempty"`
}

const (
	customerFields = ` c.cust_id, c.name, c.email, c.address,  c.city, c.stateID, c.phone, c.fax, c.contact_person, c.dealer_type,
				c.latitude,c.longitude,  c.website, c.customerID, c.isDummy, c.parentID, c.searchURL, c.eLocalURL, c.logo,c.address2,
				c.postal_code, c.mCodeID, c.salesRepID, c.APIKey, c.tier, c.showWebsite `
	stateFields            = ` IFNULL(s.state, ""), IFNULL(s.abbr, ""), IFNULL(s.countryID, "0") `
	countryFields          = ` cty.name, cty.abbr `
	dealerTypeFields       = ` IFNULL(dt.type, ""), IFNULL(dt.online, ""), IFNULL(dt.show, ""), IFNULL(dt.label, "") `
	dealerTierFields       = ` IFNULL(dtr.tier, ""), IFNULL(dtr.sort, "") `
	mapIconFields          = ` IFNULL(mi.mapicon, ""), IFNULL(mi.mapiconshadow, "") ` //joins on dealer_type usually
	mapixCodeFields        = ` IFNULL(mpx.code, ""), IFNULL(mpx.description, "") `
	salesRepFields         = ` IFNULL(sr.name, ""), IFNULL(sr.code, "") `
	customerLocationFields = ` cl.locationID, cl.name, cl.address, cl.city, cl.stateID,  cl.email, cl.phone, cl.fax,
							cl.latitude, cl.longitude, cl.cust_id, cl.contact_person, cl.isprimary, cl.postalCode, cl.ShippingDefault `
	showSiteFields = ` c.showWebsite, c.website, c.elocalurl `

	//redis
	custPrefix = "customer:"
)

func (c *Customer) GetCustomer(ctx *middleware.APIContext) (err error) {
	basicsChan := make(chan error)

	go func() {
		err := c.Basics(ctx)
		if err == nil {
			basicsChan <- c.GetUsers(ctx)
		}
		basicsChan <- err
	}()
	c.GetLocations(ctx)
	c.GetAccounts(ctx)
	err = <-basicsChan

	if err == sql.ErrNoRows {

		err = fmt.Errorf("error: %s", "failed to retrieve")
	}
	return err
}

//gets cust_id, not customerId
func (c *Customer) GetCustomerIdFromKey(ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(getCustIdFromKeyStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(ctx.DataContext.APIKey).Scan(&c.Id)
}

//redundant with Get - uses SQL joins; faster?
func (c *Customer) Basics(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(basics)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return c.ScanCustomer(stmt.QueryRow(c.Id), ctx)
}

func (c *Customer) GetLocations(ctx *middleware.APIContext) (err error) {
	redis_key := "customerLocations:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.Locations)
		return err
	}

	stmt, err := ctx.DB.Prepare(customerLocation)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(c.Id)

	for rows.Next() {
		loc, err := ScanLocation(rows)
		if err != nil {
			return err
		}
		c.Locations = append(c.Locations, *loc)
	}
	defer rows.Close()

	redis.Setex(redis_key, c.Locations, redis.CacheTimeout)

	return err
}

func (c *Customer) GetAccounts(ctx *middleware.APIContext) (err error) {
	redis_key := "CustAccount:" + strconv.Itoa(c.Id)
	data, err := redis.Get(redis_key)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &c.Accounts)
		return err
	}

	stmt, err := ctx.DB.Prepare(customerAccounts)
	if err != nil {
		return err
	}
	defer stmt.Close()

	rows, err := stmt.Query(c.Id)
	if err != nil {
		return err
	}

	for rows.Next() {
		acc, err := ScanAccount(rows)
		if err != nil {
			return err
		}

		c.Accounts = append(c.Accounts, *acc)
	}
	defer rows.Close()

	redis.Setex(redis_key, c.Accounts, redis.CacheTimeout)

	return err
}

func (c *Customer) FindCustomerIdFromCustId(ctx *middleware.APIContext) (err error) { //Jesus, really?

	stmt, err := ctx.DB.Prepare(findCustomerIdFromCustId)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(c.Id).Scan(&c.CustomerId)
}

func (c *Customer) FindCustIdFromCustomerId(ctx *middleware.APIContext) (err error) { //Jesus, really?

	stmt, err := ctx.DB.Prepare(findCustIdFromCustomerId)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(c.CustomerId).Scan(&c.Id)
}

func (c *Customer) Create(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(createCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	parentId := 0
	if c.Parent != nil {
		parentId = c.Parent.Id
	}
	res, err := stmt.Exec(
		c.Name,
		c.Email,
		c.Address,
		c.City,
		c.State.Id,
		c.Phone,
		c.Fax,
		c.ContactPerson,
		c.DealerType.Id,
		c.Latitude,
		c.Longitude,
		c.Website.String(),
		c.CustomerId,
		c.IsDummy,
		parentId,
		c.SearchUrl.String(),
		c.ELocalUrl.String(),
		c.Logo.String(),
		c.Address2,
		c.PostalCode,
		c.MapixCode.ID,
		c.SalesRepresentative.ID,
		c.ApiKey,
		c.DealerTier.Id,
		c.ShowWebsite,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	c.Id = int(id)

	for _, brandID := range c.BrandIDs {
		err = c.CreateCustomerBrand(brandID, ctx)
		if err != nil {
			return err
		}
	}
	return err
}

func (c *Customer) Update(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(updateCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()
	parentId := 0
	if c.Parent != nil {
		parentId = c.Parent.Id
	}
	_, err = stmt.Exec(
		c.Name,
		c.Email,
		c.Address,
		c.City,
		c.State.Id,
		c.Phone,
		c.Fax,
		c.ContactPerson,
		c.DealerType.Id,
		c.Latitude,
		c.Longitude,
		c.Website.String(),
		c.CustomerId,
		c.IsDummy,
		parentId,
		c.SearchUrl.String(),
		c.ELocalUrl.String(),
		c.Logo.String(),
		c.Address2,
		c.PostalCode,
		c.MapixCode.ID,
		c.SalesRepresentative.ID,
		c.ApiKey,
		c.DealerTier.Id,
		c.ShowWebsite,
		c.Id,
	)
	if err != nil {
		return err
	}
	err = c.DeleteAllCustomerBrands(ctx)
	if err != nil {
		return err
	}
	for _, brandID := range c.BrandIDs {
		err = c.CreateCustomerBrand(brandID, ctx)
	}
	go redis.Set(custPrefix+strconv.Itoa(c.Id), c)
	go redis.Delete("customerLocations:" + strconv.Itoa(c.Id))
	return nil
}

func (c *Customer) Delete(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(deleteCustomer)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Id)
	if err != nil {
		return err
	}

	err = c.DeleteAllCustomerBrands(ctx)
	if err != nil {
		return err
	}

	go redis.Delete(custPrefix + strconv.Itoa(c.Id))
	go redis.Delete("customerLocations:" + strconv.Itoa(c.Id))
	return nil
}

func (c *Customer) GetUsers(ctx *middleware.APIContext) (err error) {

	stmt, err := ctx.DB.Prepare(customerUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(c.Id)
	if err != nil {
		return err
	}
	iter := 0
	userChan := make(chan error)
	lowerKey := strings.ToLower(ctx.DataContext.APIKey)

	for res.Next() {
		var u CustomerUser
		err = res.Scan(
			&u.Id,
			&u.Name,
			&u.Email,
			&u.CustomerID,
			&u.DateAdded,
			&u.Active,
			&u.Location.Id,
			&u.Sudo,
			&u.CustID,
		)
		if err != nil {
			continue
		}

		go func(user CustomerUser) {
			if err := user.GetKeys(ctx); err == nil {
				for _, k := range user.Keys {
					if k.Key == lowerKey {
						user.Current = true
						break
					}
				}
			}
			user.Brands, err = brand.GetUserBrands(c.Id, ctx.DB)
			if err != nil {
				return
			}

			user.GetComnetAccounts(ctx)
			user.GetLocation(ctx)

			c.Users = append(c.Users, user)
			userChan <- nil

		}(u)
		iter++
	}
	defer res.Close()

	for i := 0; i < iter; i++ {
		<-userChan
	}

	return err
}

func (c *Customer) JoinUser(u CustomerUser, ctx *middleware.APIContext) error {

	stmt, err := ctx.DB.Prepare(joinUser)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(c.Id, u.Id)

	return err
}

func GetCustomerPrice(ctx *middleware.APIContext, part_id int) (price float64, err error) {

	stmt, err := ctx.DB.Prepare(customerPrice)
	if err != nil {
		return price, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(ctx.DataContext.CustomerID, part_id).Scan(&price)
	return price, err
}

func GetCustomerCartReference(ctx *middleware.APIContext, part_id int) (ref int, err error) {

	stmt, err := ctx.DB.Prepare(customerPart)
	if err != nil {
		return ref, err
	}

	err = stmt.QueryRow(ctx.DataContext.APIKey, part_id).Scan(&ref)

	return ref, err
}

//Scan Methods
func (c *Customer) ScanCustomer(res Scanner, ctx *middleware.APIContext) error {
	var err error
	var country, countryAbbr, dealerType, dealerTypeOnline, dealerTypeShow, dealerTypeLabel *string
	var dealerTier, dealerTierSort *string
	var logo, web, searchU, icon, shadow, parentId, eLocalUrl *[]byte
	var lat, lon *string
	var mapIconId, countryId *int
	var email, address, city, phone, fax, contact, address2, postal, api, state, stateAb *string
	var stateId, dType, mapixId, salesRepId, dTier, customerid *int
	var show, isdummy *bool

	err = res.Scan(
		&c.Id,
		&c.Name,
		&email,
		&address,
		&city,
		&stateId,
		&phone,
		&fax,
		&contact,
		&dType,
		&lat,
		&lon,
		&web,
		&customerid,
		&isdummy,
		&parentId,
		&searchU,
		&eLocalUrl,
		&logo,
		&address2,
		&postal,
		&mapixId,
		&salesRepId,
		&api,
		&dTier,
		&show,
		&state,
		&stateAb,
		&countryId,
		&country,
		&countryAbbr,
		&dealerType,
		&dealerTypeOnline,
		&dealerTypeShow,
		&dealerTypeLabel,
		&dealerTier,
		&dealerTierSort,
		&icon,
		&shadow,
		&c.MapixCode.Code,
		&c.MapixCode.Description,
		&c.SalesRepresentative.Name,
		&c.SalesRepresentative.Code,
	)
	if err != nil {
		return err
	}

	//get parent, if has parent
	parentChan := make(chan int)
	go func() {
		if parentId != nil {
			parentInt, err := conversions.ByteToInt(*parentId)
			if err == nil && parentInt > 0 {
				par := Customer{CustomerId: parentInt}
				if err := par.FindCustIdFromCustomerId(ctx); err != nil {
					parentChan <- 1
					return
				}

				err = par.GetCustomer(ctx)
				if err != nil {
					parentChan <- 1
					return
				}
				c.Parent = &par
			}
		}
		parentChan <- 1
	}()
	<-parentChan
	if city != nil {
		c.City = *city
	}
	if address != nil {
		c.Address = *address
	}
	if stateId != nil {
		c.State.Id = *stateId
	}
	if phone != nil {
		c.Phone = *phone
	}
	if fax != nil {
		c.Fax = *fax
	}
	if email != nil {
		c.Email = *email
	}
	if contact != nil {
		c.ContactPerson = *contact
	}
	if dType != nil {
		c.DealerType.Id = *dType
	}
	if customerid != nil {
		c.CustomerId = *customerid
	}
	if isdummy != nil {
		c.IsDummy = *isdummy
	}
	if address2 != nil {
		c.Address2 = *address2
	}
	if postal != nil {
		c.PostalCode = *postal
	}
	if api != nil {
		c.ApiKey = *api
	}
	if mapixId != nil {
		c.MapixCode.ID = *mapixId
	}
	if salesRepId != nil {
		c.SalesRepresentative.ID = *salesRepId
	}
	if dTier != nil {
		c.DealerTier.Id = *dTier
	}
	if state != nil {
		c.State.State = *state
	}
	if stateAb != nil {
		c.State.Abbreviation = *stateAb
	}

	var coun geography.Country
	if lat != nil && *lat != "" && lon != nil && *lon != "" {
		c.Latitude, _ = strconv.ParseFloat(*lat, 64)
		c.Longitude, _ = strconv.ParseFloat(*lon, 64)
	}
	if searchU != nil {
		c.SearchUrl, err = conversions.ByteToUrl(*searchU)
	}
	if eLocalUrl != nil {
		c.ELocalUrl, err = conversions.ByteToUrl(*eLocalUrl)
	}
	if logo != nil {
		c.Logo, err = conversions.ByteToUrl(*logo)
	}
	if web != nil {
		c.Website, err = conversions.ByteToUrl(*web)
	}
	if countryId != nil {
		coun.Id = *countryId
	}
	if country != nil {
		coun.Country = *country
	}
	if countryAbbr != nil {
		coun.Abbreviation = *countryAbbr
	}
	c.State.Country = &coun

	if dealerType != nil {
		c.DealerType.Type = *dealerType
	}
	if dealerTypeOnline != nil {
		c.DealerType.Online, _ = strconv.ParseBool(*dealerTypeOnline)
	}
	if dealerTypeShow != nil {
		c.DealerType.Show, _ = strconv.ParseBool(*dealerTypeShow)
	}
	if dealerTypeLabel != nil {
		c.DealerType.Label = *dealerTypeLabel
	}
	if dealerTier != nil {
		c.DealerTier.Tier = *dealerTier
	}
	if dealerTierSort != nil {
		c.DealerTier.Sort, _ = strconv.Atoi(*dealerTierSort)
	}

	if mapIconId != nil {
		c.DealerType.MapIcon.Id = *mapIconId
	}
	if icon != nil {
		c.DealerType.MapIcon.MapIcon, err = conversions.ByteToUrl(*icon)
	}
	if shadow != nil {
		c.DealerType.MapIcon.MapIconShadow, err = conversions.ByteToUrl(*shadow)
	}

	return nil
}

func ScanCustomerTableFields(res Scanner) (*Customer, error) {
	var c Customer
	var err error
	var name, email, address, address2, city, phone, fax, contactPerson, postalCode, apiKey *string
	var logo, web, searchU, parentId, eLocalUrl *[]byte
	var lat, lon *string
	var stateId, dtypeId, dtierId, custID, mapixCodeID, salesRepID *int
	var isDummy, showWebsite *bool

	err = res.Scan(
		&c.Id,
		&name,
		&email,
		&address,
		&city,
		&stateId,
		&phone,
		&fax,
		&contactPerson,
		&dtypeId,
		&lat,
		&lon,
		&web,
		&custID,
		&isDummy,
		&parentId,
		&searchU,
		&eLocalUrl,
		&logo,
		&address2,
		&postalCode,
		&mapixCodeID,
		&salesRepID,
		&apiKey,
		&dtierId,
		&showWebsite,
	)
	if err != nil {
		return &c, err
	}

	if name != nil {
		c.Name = *name
	}
	if address != nil {
		c.Address = *address
	}
	if address2 != nil {
		c.Address2 = *address2
	}
	if city != nil {
		c.City = *city
	}
	if email != nil {
		c.Email = *email
	}
	if phone != nil {
		c.Phone = *phone
	}
	if fax != nil {
		c.Fax = *fax
	}
	if contactPerson != nil {
		c.ContactPerson = *contactPerson
	}
	if lat != nil && *lat != "" && lon != nil && *lon != "" {
		c.Latitude, _ = strconv.ParseFloat(*lat, 64)
		c.Longitude, _ = strconv.ParseFloat(*lon, 64)
	}
	if searchU != nil {
		c.SearchUrl, err = conversions.ByteToUrl(*searchU)
	}
	if eLocalUrl != nil {
		c.ELocalUrl, err = conversions.ByteToUrl(*eLocalUrl)
	}
	if logo != nil {
		c.Logo, err = conversions.ByteToUrl(*logo)
	}
	if web != nil {
		c.Website, err = conversions.ByteToUrl(*web)
	}
	if custID != nil {
		c.CustomerId = *custID
	}
	if isDummy != nil {
		c.IsDummy = *isDummy
	}
	if postalCode != nil {
		c.PostalCode = *postalCode
	}
	if mapixCodeID != nil {
		c.MapixCode.ID = *mapixCodeID
	}
	if salesRepID != nil {
		c.SalesRepresentative.ID = *salesRepID
	}
	if apiKey != nil {
		c.ApiKey = *apiKey
	}
	if showWebsite != nil {
		c.ShowWebsite = *showWebsite
	}
	if stateId != nil {
		c.State.Id = *stateId
	}
	if dtypeId != nil {
		c.DealerType.Id = *dtypeId
	}
	if dtierId != nil {
		c.DealerTier.Id = *dtierId
	}
	return &c, err
}

func ScanAccount(res Scanner) (*Account, error) {
	var a Account
	var err error

	var accID *int
	var accountNumber *string
	var cust_id *int
	var typeID *int
	var typeText *string
	var comnetURL *[]byte
	var freightLimit *float64
	var defaultWare *int

	err = res.Scan(
		&accID,
		&accountNumber,
		&cust_id,
		&typeID,
		&freightLimit,
		&typeText,
		&comnetURL,
		&defaultWare,
	)
	if err != nil {
		return &a, err
	}
	if accID != nil {
		a.ID = *accID
	}
	if accountNumber != nil {
		a.AccountNumber = *accountNumber
	}
	if cust_id != nil {
		a.Cust_id = *cust_id
	}
	if typeID != nil {
		a.TypeID = *typeID
		a.Type.ID = a.TypeID
	}
	if freightLimit != nil {
		a.FreightLimit = *freightLimit
	}
	if typeText != nil {
		a.Type.Title = *typeText
	}
	if comnetURL != nil {
		a.Type.ComnetURL, err = conversions.ByteToUrl(*comnetURL)
	}
	if defaultWare != nil {
		a.DefaultWarehouseID = *defaultWare
	}

	return &a, err
}

func ScanLocation(res Scanner) (*CustomerLocation, error) {
	var l CustomerLocation
	var err error
	var name, email, address, city, phone, fax, contactPerson, state, stateAbbr, country, countryAbbr, postalCode *string
	var lat, lon *float64
	var custId, stateId, countryId *int
	var isPrimary, shippingDefault *bool
	var coun geography.Country

	err = res.Scan(
		&l.Id,
		&name,
		&address,
		&city,
		&stateId,
		&email,
		&phone,
		&fax,
		&lat,
		&lon,
		&custId,
		&contactPerson,
		&isPrimary,
		&postalCode,
		&shippingDefault,
		&state,
		&stateAbbr,
		&countryId,
		&country,
		&countryAbbr,
	)
	if err != nil {
		return &l, err
	}
	if name != nil {
		l.Name = *name
	}
	if email != nil {
		l.Email = *email
	}
	if address != nil {
		l.Address = *address
	}
	if city != nil {
		l.City = *city
	}
	if postalCode != nil {
		l.PostalCode = *postalCode
	}
	if phone != nil {
		l.Phone = *phone
	}
	if fax != nil {
		l.Fax = *fax
	}
	if lat != nil {
		l.Coordinates.Latitude = *lat
	}
	if lon != nil {
		l.Coordinates.Longitude = *lon
	}
	if custId != nil {
		l.CustomerId = *custId
	}
	if contactPerson != nil {
		l.ContactPerson = *contactPerson
	}
	if isPrimary != nil {
		l.IsPrimary = *isPrimary
	}
	if shippingDefault != nil {
		l.ShippingDefault = *shippingDefault
	}
	if stateId != nil {
		l.State.Id = *stateId
	}
	if state != nil {
		l.State.State = *state
	}
	if stateAbbr != nil {
		l.State.Abbreviation = *stateAbbr
	}
	if countryId != nil {
		coun.Id = *countryId
	}
	if country != nil {
		coun.Country = *country
	}
	if countryAbbr != nil {
		coun.Abbreviation = *countryAbbr
	}
	l.State.Country = &coun
	return &l, err
}

func ScanDealerLocation(res *sql.Rows, count int) (*DealerLocation, error) {
	var l DealerLocation
	var err error

	vals := make([]sql.RawBytes, count)
	scanArgs := make([]interface{}, count)
	for i := range vals {
		scanArgs[i] = &vals[i]
	}

	err = res.Scan(scanArgs...)
	if err != nil {
		return &l, err
	}

	if vals[0] != nil {
		l.CustomerLocation.Id, _ = conversions.ByteToInt(vals[0])
	}

	if vals[1] != nil {
		l.Name, _ = conversions.ByteToString(vals[1])
	}
	if vals[2] != nil {
		l.Address, _ = conversions.ByteToString(vals[2])
	}
	if vals[3] != nil {
		l.City, _ = conversions.ByteToString(vals[3])
	}
	if vals[4] != nil {
		l.State.Id, _ = conversions.ByteToInt(vals[4])
	}
	if vals[5] != nil {
		l.Email, _ = conversions.ByteToString(vals[5])
	}
	if vals[6] != nil {
		l.Phone, _ = conversions.ByteToString(vals[6])
	}
	if vals[7] != nil {
		l.Fax, _ = conversions.ByteToString(vals[7])
	}
	if vals[8] != nil {
		l.Coordinates.Latitude, _ = conversions.ByteToFloat(vals[8])
	}
	if vals[9] != nil {
		l.Coordinates.Latitude, _ = conversions.ByteToFloat(vals[9])
	}
	if vals[10] != nil {
		l.CustomerId, _ = conversions.ByteToInt(vals[10])
	}
	if vals[11] != nil {
		l.ContactPerson, _ = conversions.ByteToString(vals[11])
	}
	if vals[12] != nil {
		l.IsPrimary, _ = conversions.ParseBool(vals[12])
	}
	if vals[13] != nil {
		l.PostalCode, _ = conversions.ByteToString(vals[13])
	}
	if vals[14] != nil {
		l.ShippingDefault, _ = conversions.ParseBool(vals[14])
	}
	if vals[15] != nil {
		l.State.State, _ = conversions.ByteToString(vals[15])
	}
	if vals[16] != nil {
		l.State.Abbreviation, _ = conversions.ByteToString(vals[16])
	}
	if vals[17] != nil {
		if l.State.Country == nil {
			l.State.Country = &geography.Country{}
		}
		l.State.Country.Id, _ = conversions.ByteToInt(vals[17])
	}
	if vals[18] != nil {
		if l.State.Country == nil {
			l.State.Country = &geography.Country{}
		}
		l.State.Country.Country, _ = conversions.ByteToString(vals[18])
	}
	if vals[19] != nil {
		l.State.Country.Abbreviation, _ = conversions.ByteToString(vals[19])
	}
	if vals[20] != nil {
		l.DealerType.Type, _ = conversions.ByteToString(vals[20])
	}
	if vals[21] != nil {
		l.DealerType.Online, _ = conversions.ParseBool(vals[21])
	}
	if vals[22] != nil {
		l.DealerType.Show, _ = conversions.ParseBool(vals[22])
	}
	if vals[23] != nil {
		l.DealerType.Label, _ = conversions.ByteToString(vals[23])
	}
	if vals[24] != nil {
		l.DealerTier.Tier, _ = conversions.ByteToString(vals[24])
	}
	if vals[25] != nil {
		l.DealerTier.Sort, _ = conversions.ByteToInt(vals[25])
	}
	if vals[26] != nil {
		l.DealerType.MapIcon.MapIcon, _ = conversions.ByteToUrl(vals[26])
	}
	if vals[27] != nil {
		l.DealerType.MapIcon.MapIconShadow, _ = conversions.ByteToUrl(vals[27])
	}
	if vals[28] != nil {
		l.MapixCode.Code, _ = conversions.ByteToString(vals[28])
	}
	if vals[29] != nil {
		l.MapixCode.Description, _ = conversions.ByteToString(vals[29])
	}
	if vals[30] != nil {
		l.SalesRepresentative.Name, _ = conversions.ByteToString(vals[30])
	}
	if vals[31] != nil {
		l.SalesRepresentative.Code, _ = conversions.ByteToString(vals[31])
	}
	if vals[32] != nil {
		l.ShowWebSite, _ = conversions.ParseBool(vals[32])
	}
	if vals[33] != nil {
		l.Website, _ = conversions.ByteToUrl(vals[33])
	}
	if vals[34] != nil {
		l.ELocalUrl, _ = conversions.ByteToUrl(vals[34])
	}

	if len(vals) == 36 && vals[35] != nil {
		l.Distance, _ = conversions.ByteToFloat(vals[35])
	}

	return &l, err
}

var (
	getCustomer              = `select ` + customerFields + ` from Customer as c where c.cust_id = ? `
	getCustomerIdFromKeyStmt = `select c.customerID from Customer as c
                                join CustomerUser as cu on c.cust_id = cu.cust_ID
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`

	getCustIdFromKeyStmt = `select cu.cust_ID from CustomerUser as cu
                                join ApiKey as ak on cu.id = ak.user_id
                                where ak.api_key = ? limit 1`
	//Old
	findCustomerIdFromCustId = `select customerID from Customer where cust_id = ? limit 1`
	findCustIdFromCustomerId = `select cust_id from Customer where customerID = ? limit 1`
	basics                   = `select ` + customerFields + `, ` + stateFields + `, ` + countryFields + `, ` + dealerTypeFields + `, ` + dealerTierFields + `, ` + mapIconFields + `, ` + mapixCodeFields + `, ` + salesRepFields + `
			from Customer as c
				left join States as s on c.stateID = s.stateID
				left join Country as cty on s.countryID = cty.countryID
				left join DealerTypes as dt on c.dealer_type = dt.dealer_type
				left join MapIcons as mi on dt.dealer_type = mi.dealer_type
				left join DealerTiers as dtr on c.tier = dtr.ID
				left join MapixCode as mpx on c.mCodeID = mpx.mCodeID
				left join SalesRepresentative as sr on c.salesRepID = sr.salesRepID
				where c.cust_id = ? `
	//affects CL methods
	customerLocation = `select ` + customerLocationFields + `, ` + stateFields + `, ` + countryFields + `
							from CustomerLocations as cl
							join Customer as c on cl.cust_id = c.cust_id
	 						left join States as s on cl.stateID = s.stateID
	 						left join Country as cty on s.countryID = cty.countryID
							where c.cust_id = ?`
	customerAccounts = `select act.id, act.accountNumber, act.cust_id, act.type_id, act.freightLimit, acty.type, acty.comnet_url, act.defaultWarehouseId from Accounts as act
							Join AccountTypes as acty on acty.id = act.type_id
							where act.cust_id = ?`

	customerUser = `select cu.id, cu.name, cu.email, cu.customerID, cu.date_added, cu.active,cu.locationID, cu.isSudo, cu.cust_ID from CustomerUser as cu
						join Customer as c on cu.cust_ID = c.cust_id
						where c.cust_id = ?
						&& cu.active = 1`
	customerPrice = `select distinct cp.price from
						 CustomerPricing cp
						 where cp.cust_ID =  ?
						 and cp.partID = ?`

	customerPart = `select distinct ci.custPartID from ApiKey as ak
						join CustomerUser cu on ak.user_id = cu.id
						join Customer c on cu.cust_ID = c.cust_id
						join CartIntegration ci on c.cust_ID = ci.custID
						where ak.api_key = ?
						and ci.partID = ?`

	//customer Crud
	createCustomer = `insert into Customer (name, email, address,  city, stateID, phone, fax, contact_person, dealer_type, latitude,longitude,  website, customerID, isDummy, parentID, searchURL,
					eLocalURL, logo,address2, postal_code, mCodeID, salesRepID, APIKey, tier, showWebsite) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`

	updateCustomer = `update Customer set name = ?, email = ?, address = ?, city = ?, stateID = ?, phone = ?, fax = ?, contact_person = ?, dealer_type = ?, latitude = ?, longitude = ?,  website = ?, customerID = ?,
					isDummy = ?, parentID = ?, searchURL = ?, eLocalURL = ?, logo = ?, address2 = ?, postal_code = ?, mCodeID = ?, salesRepID = ?, APIKey = ?, tier = ?, showWebsite = ? where cust_id = ?`
	deleteCustomer   = `delete from Customer where cust_id = ?`
	joinUser         = `update CustomerUser set cust_ID = ? where id = ?`
	createDealerType = `insert into DealerTypes (type, online, label) values(?,?,?)`
	deleteDealerType = `delete from DealerTypes where dealer_type= ?`
)
