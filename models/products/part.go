package products

import (
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

// Part ...
type Part struct {
	Identifier        bson.ObjectId        `bson:"_id" json:"-" xml:"-"`
	ID                int                  `json:"id" xml:"id,attr" bson:"id"`
	PartNumber        string               `bson:"part_number" json:"part_number" xml:"part_number,attr"`
	Brand             brand.Brand          `json:"brand" xml:"brand,attr" bson:"brand"`
	Status            int                  `json:"status" xml:"status,attr" bson:"status"`
	PriceCode         int                  `json:"price_code" xml:"price_code,attr" bson:"price_code"`
	RelatedCount      int                  `json:"related_count" xml:"related_count,attr" bson:"related_count"`
	AverageReview     float64              `json:"average_review" xml:"average_review,attr" bson:"average_review"`
	DateModified      time.Time            `json:"date_modified" xml:"date_modified,attr" bson:"date_modified"`
	DateAdded         time.Time            `json:"date_added" xml:"date_added,attr" bson:"date_added"`
	ShortDesc         string               `json:"short_description" xml:"short_description,attr" bson:"short_description"`
	InstallSheet      *url.URL             `json:"install_sheet" xml:"install_sheet" bson:"install_sheet"`
	Attributes        []Attribute          `json:"attributes" xml:"attributes" bson:"attributes"`
	AcesVehicles      []AcesVehicle        `bson:"aces_vehicles" json:"aces_vehicles" xml:"aces_vehicles"`
	VehicleAttributes []string             `json:"vehicle_atttributes" xml:"vehicle_attributes" bson:"vehicle_attributes"`
	Vehicles          []VehicleApplication `json:"vehicle_applications,omitempty" xml:"vehicle_applications,omitempty" bson:"vehicle_applications"`
	Content           []Content            `json:"content" xml:"content" bson:"content"`
	Pricing           []Price              `json:"pricing" xml:"pricing" bson:"pricing"`
	Reviews           []Review             `json:"reviews" xml:"reviews" bson:"reviews"`
	Images            []Image              `json:"images" xml:"images" bson:"images"`
	Related           []int                `json:"related" xml:"related" bson:"related" bson:"related"`
	Categories        []Category           `json:"categories" xml:"categories" bson:"categories"`
	Videos            []video.Video        `json:"videos" xml:"videos" bson:"videos"`
	Packages          []Package            `json:"packages" xml:"packages" bson:"packages"`
	Customer          CustomerPart         `json:"customer,omitempty" xml:"customer,omitempty" bson:"v"`
	Class             Class                `json:"class,omitempty" xml:"class,omitempty" bson:"class"`
	Featured          bool                 `json:"featured,omitempty" xml:"featured,omitempty" bson:"featured"`
	AcesPartTypeID    int                  `json:"acesPartTypeId,omitempty" xml:"acesPartTypeId,omitempty" bson:"acesPartTypeId"`
	Inventory         PartInventory        `json:"inventory,omitempty" xml:"inventory,omitempty" bson:"inventory"`
	UPC               string               `json:"upc,omitempty" xml:"upc,omitempty" bson:"upc"`
}

// CustomerPart Holds customer specific data.
type CustomerPart struct {
	Price         float64 `json:"price" xml:"price,attr"`
	CartReference int     `json:"cart_reference" xml:"cart_reference,attr"`
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

// VehicleApplication ...
type VehicleApplication struct {
	Year        string `bson:"year" json:"year" xml:"year,attr"`
	Make        string `bson:"make" json:"make" xml:"make,attr"`
	Model       string `bson:"model" json:"model" xml:"model,attr"`
	Style       string `bson:"style" json:"style" xml:"style,attr"`
	Exposed     string `bson:"exposed" json:"exposed" xml:"exposed"`
	Drilling    string `bson:"drilling" json:"drilling" xml:"drilling"`
	InstallTime string `bson:"install_time" json:"install_time" xml:"install_time"`
}

// Identifiers ...
func Identifiers(ctx *middleware.APIContext, brandID int) ([]string, error) {
	var parts []string
	brands := []int{brandID}

	if brandID == 0 && ctx.DataContext.BrandID != 0 {
		brands = []int{ctx.DataContext.BrandID}
	} else if brandID == 0 && ctx.DataContext.BrandID == 0 {
		brands = ctx.DataContext.BrandArray
	}

	qry := bson.M{
		"brand.id": bson.M{
			"$in": brands,
		},
		"status": bson.M{
			"$in": []int{800, 900},
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
func All(page, count int, ctx *middleware.APIContext) ([]Part, error) {
	var parts []Part
	brands := []int{ctx.Brand}
	if ctx.Brand == 0 {
		brands = ctx.DataContext.BrandArray
	}

	qry := bson.M{
		"brand.id": bson.M{
			"$in": brands,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Sort("id:1").Skip(page * count).Limit(count).All(&parts)

	return parts, err
}

// Featured Returns `count` featured products.
func Featured(ctx *middleware.APIContext, count, brandID int) ([]Part, error) {
	var parts []Part
	brands := []int{brandID}

	if brandID == 0 && ctx.DataContext.BrandID != 0 {
		brands = []int{ctx.DataContext.BrandID}
	} else if brandID == 0 && ctx.DataContext.BrandID == 0 {
		brands = ctx.DataContext.BrandArray
	}

	qry := bson.M{
		"featured": true,
		"brand.id": bson.M{
			"$in": brands,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Sort("id:1").Limit(count).All(&parts)

	return parts, err
}

// Latest Returns `count` latest products.
func Latest(ctx *middleware.APIContext, count, brandID int) ([]Part, error) {
	var parts []Part
	brands := []int{brandID}

	if brandID == 0 && ctx.DataContext.BrandID != 0 {
		brands = []int{ctx.DataContext.BrandID}
	} else if brandID == 0 && ctx.DataContext.BrandID == 0 {
		brands = ctx.DataContext.BrandArray
	}

	qry := bson.M{
		"brand.id": bson.M{
			"$in": brands,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).Sort("-date_added").Limit(count).All(&parts)
	return parts, err
}

// GetRelated Returns all the products that are related to this one.
func (p *Part) GetRelated(ctx *middleware.APIContext, brandID int) ([]Part, error) {
	var parts []Part
	brands := []int{brandID}

	if brandID == 0 && ctx.DataContext.BrandID != 0 {
		brands = []int{ctx.DataContext.BrandID}
	} else if brandID == 0 && ctx.DataContext.BrandID == 0 {
		brands = ctx.DataContext.BrandArray
	}

	query := bson.M{
		"id": bson.M{
			"$in": p.Related,
		},
		"brand.id": bson.M{
			"$in": brands,
		},
	}

	err := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(query).Sort("id:1").All(&parts)

	return parts, err
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
func (p *Part) Get(ctx *middleware.APIContext, brandID int) (err error) {
	brands := []int{brandID}

	if brandID == 0 && ctx.DataContext.BrandID != 0 {
		brands = []int{ctx.DataContext.BrandID}
	} else if brandID == 0 && ctx.DataContext.BrandID == 0 {
		brands = ctx.DataContext.BrandArray
	}

	pattern := bson.RegEx{
		Pattern: "^" + p.PartNumber + "$",
		Options: "i",
	}

	qry := bson.M{
		"part_number": pattern,
		"brand.id":    bson.M{"$in": brands},
	}

	return ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Find(qry).One(&p)
}
