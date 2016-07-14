package products

import (
	"fmt"
	"sort"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/customer/content"
	"github.com/curt-labs/API/models/video"

	"gopkg.in/mgo.v2/bson"

	"net/url"
	"strings"
	"time"
)

const (
	// MaxPartCount defines the maximum numbers of products that can be returned
	// for a paginated query.
	MaxPartCount = 250
)

// Part ...
type Part struct {
	Identifier    bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	ID            int           `json:"id" xml:"id,attr" bson:"id"`
	SKU           string        `bson:"part_number" json:"sku" xml:"sku,attr"`
	Brand         brand.Brand   `json:"brand" xml:"brand,attr" bson:"brand"`
	Status        int           `json:"status" xml:"status,attr" bson:"status"`
	PriceCode     int           `json:"price_code" xml:"price_code,attr" bson:"price_code"`
	RelatedCount  int           `json:"related_count" xml:"related_count,attr" bson:"related_count"`
	AverageReview float64       `json:"average_review" xml:"average_review,attr" bson:"average_review"`
	DateModified  time.Time     `json:"date_modified" xml:"date_modified,attr" bson:"date_modified"`
	DateAdded     time.Time     `json:"date_added" xml:"date_added,attr" bson:"date_added"`
	ShortDesc     string        `json:"short_description" xml:"short_description,attr" bson:"short_description"`
	InstallSheet  *url.URL      `json:"install_sheet" xml:"install_sheet" bson:"install_sheet"`
	Attributes    []Attribute   `json:"attributes" xml:"attributes" bson:"attributes"`

	// TODO: This needs to be re-evaluated for the integrated new web.
	// AcesVehicles      []AcesVehicle        `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`

	// TODO: Is this really needed? All of this data is available on the
	// `Vehicles` array. It's more of a "nice to have", worth the bandwidth?
	// VehicleAttributes []string             `json:"vehicle_atttributes" xml:"vehicle_attributes" bson:"vehicle_attributes"`

	Vehicles       []VehicleApplication `json:"vehicle_applications,omitempty" xml:"vehicle_applications,omitempty" bson:"vehicle_applications"`
	Content        []Content            `json:"content" xml:"content" bson:"content"`
	Pricing        []Price              `json:"pricing" xml:"pricing" bson:"pricing"`
	Reviews        []Review             `json:"reviews" xml:"reviews" bson:"reviews"`
	Images         []Image              `json:"images" xml:"images" bson:"images"`
	Related        []int                `json:"related" xml:"related" bson:"related" bson:"related"`
	Categories     []Category           `json:"categories" xml:"categories" bson:"categories"`
	Videos         []video.Video        `json:"videos" xml:"videos" bson:"videos"`
	Packages       []Package            `json:"packages" xml:"packages" bson:"packages"`
	Customer       CustomerPart         `json:"customer,omitempty" xml:"customer,omitempty" bson:"v"`
	Class          Class                `json:"class,omitempty" xml:"class,omitempty" bson:"class"`
	Featured       bool                 `json:"featured,omitempty" xml:"featured,omitempty" bson:"featured"`
	AcesPartTypeID int                  `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty" bson:"acesPartTypeId"`
	Inventory      PartInventory        `json:"inventory,omitempty" xml:"inventory,omitempty" bson:"inventory"`
	UPC            string               `json:"upc,omitempty" xml:"upc,omitempty" bson:"upc"`
}

// CustomerPart Holds customer specific data.
type CustomerPart struct {
	Price           float64                   `json:"price" xml:"price,attr"`
	CartReference   int                       `json:"cart_reference" xml:"cart_reference,attr"`
	CustomerContent []customerContent.Content `json:"content" xml:"content"`
}

// PaginatedProductListing ...
type PaginatedProductListing struct {
	Parts         []Part `json:"parts" xml:"parts"`
	TotalItems    int    `json:"total_items" xml:"total_items,attr"`
	ReturnedCount int    `json:"returned_count" xml:"returned_count,attr"`
	Page          int    `json:"page" xml:"page,attr"`
	PerPage       int    `json:"per_page" xml:"per_page,attr"`
	TotalPages    int    `json:"total_pages" xml:"total_pages,attr"`
}

// Class ...
type Class struct {
	ID    int    `json:"id,omitempty" xml:"id,omitempty" bson:"id"`
	Name  string `json:"name,omitempty" xml:"name,omitempty" bson:"name"`
	Image string `json:"image,omitempty" xml:"image,omitempty" bson:"image"`
}

// Identifiers ...
func Identifiers(ctx *middleware.APIContext) ([]string, error) {
	var parts []string

	qry := bson.M{
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
		"status": bson.M{
			"$in": ctx.Statuses,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Distinct("part_number", &parts)
	if err != nil {
		return []string{}, err
	}

	sort.Strings(parts)

	return parts, nil
}

// All Retrieves all parts for the given brands
func All(ctx *middleware.APIContext, page, count int) ([]Part, error) {
	var parts []Part

	if count > MaxPartCount {
		count = MaxPartCount
	}

	qry := bson.M{
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}

	col := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	err := col.Find(qry).Sort("id:1").Skip(page * count).Limit(count).All(&parts)

	return parts, err
}

// Featured Returns `count` featured products.
func Featured(ctx *middleware.APIContext, count int) ([]Part, error) {
	var parts []Part

	qry := bson.M{
		"featured": true,
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Sort("id:1").Limit(count).All(&parts)

	return parts, err
}

// Latest Returns `count` latest products.
func Latest(ctx *middleware.APIContext, count int) ([]Part, error) {
	var parts []Part

	qry := bson.M{
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Sort("-date_added").Limit(count).All(&parts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve latest from DB: %s\n", err.Error())
	}

	return parts, nil
}

// GetRelated Returns all the products that are related to this one.
func (p *Part) GetRelated(ctx *middleware.APIContext) ([]Part, error) {
	var parts []Part
	col := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	if len(p.Related) == 0 {
		// part hasn't retrieved it's related array, so we'll do it
		err := col.Find(bson.M{"part_number": p.SKU}).Select(bson.M{"related": 1, "_id": 0}).All(&p.Related)
		if err != nil {
			return nil, fmt.Errorf("failed to lookup related identifiers: %s\n", err.Error())
		}
	}
	query := bson.M{
		"id": bson.M{
			"$in": p.Related,
		},
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}

	err := col.Find(query).Sort("id:1").All(&parts)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve related parts from DB: %s\n", err.Error())
	}

	return parts, nil
}

// BindCustomer Binds customer specific data to the product.
func (p *Part) BindCustomer(ctx *middleware.APIContext) CustomerPart {
	var price float64
	var ref int

	priceChan := make(chan int)
	refChan := make(chan int)

	go func() {
		price, _ = customer.GetCustomerPrice(ctx.DB, p.ID)
		priceChan <- 1
	}()

	go func() {
		ref, _ = customer.GetCustomerCartReference(ctx.DB, p.ID)
		refChan <- 1
	}()

	content, _ := customerContent.GetPartContent(ctx.DB, p.ID)
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

	<-priceChan
	<-refChan

	return CustomerPart{
		Price:         price,
		CartReference: ref,
	}
}

// Get Retrieves product data.
func (p *Part) Get(ctx *middleware.APIContext) (err error) {

	pattern := bson.RegEx{
		Pattern: "^" + p.SKU + "$",
		Options: "i",
	}

	qry := bson.M{
		"part_number": pattern,
		"brand.id":    bson.M{"$in": ctx.DataContext.BrandArray},
	}

	return ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).One(&p)
}

// GetAttributes ...
func GetAttributes(ctx *middleware.APIContext, sku string) ([]Attribute, error) {

	pattern := bson.RegEx{
		Pattern: "^" + sku + "$",
		Options: "i",
	}

	qry := bson.M{
		"part_number": pattern,
		"brand.id":    bson.M{"$in": ctx.DataContext.BrandArray},
	}

	col := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)
	var pa []Attribute

	err := col.Find(qry).Select(bson.M{"part_attributes": 1, "_id": 0}).One(&pa)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product attributes: %s\n", err.Error())
	}

	return pa, nil
}

// GetVehicles ...
func GetVehicles(ctx *middleware.APIContext, sku string) ([]VehicleApplication, error) {

	pattern := bson.RegEx{
		Pattern: "^" + sku + "$",
		Options: "i",
	}

	qry := bson.M{
		"part_number": pattern,
		"brand.id":    bson.M{"$in": ctx.DataContext.BrandArray},
	}

	col := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)
	var va []VehicleApplication

	err := col.Find(qry).Select(bson.M{"vehicle_applications": 1, "_id": 0}).One(&va)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vehicle applications: %s\n", err.Error())
	}

	return va, nil
}
