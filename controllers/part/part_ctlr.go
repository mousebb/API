package partCtlr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/rest"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/vehicle"
	"github.com/ninnemana/analytics-go"
)

func track(endpoint string, params map[string]string, r *http.Request) {
	client := analytics.New("sud7rjoq3o")
	client.FlushAfter = 30 * time.Second
	client.FlushAt = 25

	js, err := json.Marshal(params)
	if err != nil {
		return
	}

	client.Track(map[string]interface{}{
		"title":    "Part Endpoint",
		"url":      r.URL.String(),
		"path":     r.URL.Path,
		"referrer": r.URL.RequestURI(),
		"params":   js,
	})
}

// Identifiers Returns a slice of distinct part numbers.
func Identifiers(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var b int
	if r.URL.Query().Get("brand") != "" {
		b, _ = strconv.Atoi(r.URL.Query().Get("brand"))
	}

	return products.Identifiers(ctx, b)
}

// All Returns a slice of all Part.
func All(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	page := 0
	count := 10
	qs := r.URL.Query()

	if qs.Get("page") != "" {
		if pg, err := strconv.Atoi(qs.Get("page")); err == nil {
			if pg == 0 {
				pg = 1
			}
			page = pg - 1
		}
	}
	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 500 {
				return nil, fmt.Errorf("maximum request size is 500, you requested: %d", ct)
			}
			count = ct
		}
	}

	return products.All(page, count, ctx)
}

// Featured Returns a given amount of featured Part.
func Featured(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				return nil, fmt.Errorf("maximum request size is 50, you requested: %d", ct)
			}
			count = ct
		}
	}
	var brand int
	var err error
	brandStr := qs.Get("brand")
	if brandStr != "" {
		brand, err = strconv.Atoi(brandStr)
		if err != nil {
			return nil, err
		}
	}

	return products.Featured(ctx, count, brand)
}

// Latest Returns the latest slice of Part.
func Latest(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				return nil, fmt.Errorf("maximum request size is 50, you requested: %d", ct)
			}
			count = ct
		}
	}
	var brand int
	var err error
	brandStr := qs.Get("brand")
	if brandStr != "" {
		brand, err = strconv.Atoi(brandStr)
		if err != nil {
			return nil, err
		}
	}

	return products.Latest(ctx, count, brand)
}

// Get Retrieves a Part using the :part segment of the URL pattern.
// If it's an ARIES product it binds the MongoDB vehicle data.
// TODO: We should add logic for the CURT year/make/model/style.
func Get(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)
	if err != nil {
		return nil, err
	}

	// QUESTION: Is this the right place for the vehicle logic? Should
	// it be handled in the fanner process instead? I would think this data
	// would have already been applied before indexing.

	//TODO - remove when curt & aries vehicle application data are in sync
	if p.Brand.ID == 3 {
		mgoVehicles, err := vehicle.ReverseMongoLookup(p.ID, ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range mgoVehicles {
			vehicleApplication := products.VehicleApplication{
				Year:  v.Year,
				Make:  v.Make,
				Model: v.Model,
				Style: v.Style,
			}
			p.Vehicles = append(p.Vehicles, vehicleApplication)
		}
	} //END TODO

	return p, nil
}

// GetRelated Retrieves the related Part to a given Part.
func GetRelated(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	return p.GetRelated(ctx, 0)
}

// GetWithVehicle Gets a Part with attributes relative to the fitment
// to a Vehicle.
func GetWithVehicle(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	err = errors.New("Not Implemented")
	return nil, err
	// qs := r.URL.Query()
	// partID, err := strconv.Atoi(params["part"])
	// if err != nil {
	// 	http.Error(w, "Invalid part number", http.StatusInternalServerError)
	// 	return ""
	// }
	// key := qs.Get("key")
	// year, err := strconv.ParseFloat(params["year"], 64)
	// if err != nil {
	// 	http.Redirect(w, r, "/part/"+params["part"]+"?key="+key, http.StatusFound)
	// 	return ""
	// }
	// vMake := params["make"]
	// model := params["model"]
	// submodel := params["submodel"]
	// config_vals := strings.Split(strings.TrimSpace(params["config"]), "/")

	// vehicle := Vehicle{
	// 	Year:          year,
	// 	Make:          vMake,
	// 	Model:         model,
	// 	Submodel:      submodel,
	// 	Configuration: config_vals,
	// }

	// p := products.Part{
	// 	ID: partID,
	// }

	// err = products.GetWithVehicle(&vehicle, key)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return ""
	// }

	// return encoding.Must(enc.Encode(part))
}

// Vehicles Returns the vehicles that fit a given Part.
func Vehicles(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)
	if err != nil {
		return err, nil
	}

	return vehicle.ReverseLookup(ctx, p.ID)
}

//Redundant
func Images(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	return p.Images, err
}

//Redundant
func Attributes(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	return p.Attributes, err
}

//Redundant
func GetContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	return p.Content, err
}

//Redundant
func Packaging(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	return p.Packages, err
}

//Redundant
func ActiveApprovedReviews(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)
	if err != nil {
		return nil, err
	}

	var revs []products.Review
	for _, rev := range p.Reviews {
		if rev.Active == true && rev.Approved == true {
			revs = append(revs, rev)
		}
	}

	return revs, nil
}

func Videos(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)
	if err != nil {
		return nil, err
	}

	return p.Videos, err
}

//Sort of Redundant
func InstallSheet(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.Split(ctx.Params.ByName("part"), ".")[0])
	if err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return
	}
	p := products.Part{
		ID: id,
	}

	err = p.Get(ctx, 0)
	if err != nil {
		apierror.GenerateError("Trouble getting part", err, rw, r)
		return
	}
	var text string
	for _, content := range p.Content {
		if content.ContentType.Type == "installationSheet" {
			text = content.Text
		}
	}
	if text == "" {
		apierror.GenerateError("No Installation Sheet", err, rw, r)
		return
	}

	data, err := rest.GetPDF(text, r)
	if err != nil {
		apierror.GenerateError("Error getting PDF", err, rw, r)
		return
	}

	rw.Header().Set("Content-Length", strconv.Itoa(len(data)))
	rw.Header().Set("Content-Type", "application/pdf")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin")
	rw.Write(data)
}

// Categories Returns product categories.
func Categories(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	return p.Categories, err
}

func Prices(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		PartNumber: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx, 0)

	custChan := make(chan products.Price)

	go func() {
		price, custErr := customer.GetCustomerPrice(ctx, p.ID)
		if custErr != nil {
			err = custErr
		}
		custChan <- products.Price{0, 0, "Customer", price, false, time.Now()}
	}()

	err = p.Get(ctx, 0)
	if err != nil {
		return nil, err
	}

	p.Pricing = append(p.Pricing, <-custChan)

	return p.Pricing, nil
}
