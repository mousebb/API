package vehicle

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/curt-labs/API/helpers/apifilter"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/vehicle"
)

var (
	ignoredFormParams = []string{"key"}
)

// Finds further configuration options and parts that match
// the given configuration. Doesn't start looking for parts
// until the model is provided.
func Query(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var l products.Lookup
	var page int
	var count int

	qs := r.URL.Query()

	page, _ = strconv.Atoi(qs.Get("page"))
	count, _ = strconv.Atoi(qs.Get("count"))
	qs.Del("page")
	qs.Del("count")

	l.Vehicle = LoadVehicle(r)

	l.Brands = ctx.DataContext.BrandArray

	if qs.Get("key") != "" {
		l.CustomerKey = qs.Get("key")
	} else if r.FormValue("key") != "" {
		l.CustomerKey = r.FormValue("key")
		delete(r.Form, "key")
	} else {
		l.CustomerKey = r.Header.Get("key")
	}

	if l.Vehicle.Base.Year == 0 { // Get Years
		if err := l.GetYears(ctx); err != nil {
			return nil, errors.New("Trouble getting years for vehicle lookup")
		}
	} else if l.Vehicle.Base.Make == "" { // Get Makes
		if err := l.GetMakes(ctx); err != nil {
			return nil, errors.New("Trouble getting makes for vehicle lookup")
		}
	} else if l.Vehicle.Base.Model == "" { // Get Models
		if err := l.GetModels(ctx); err != nil {
			return nil, errors.New("Trouble getting models for vehicle lookup")
		}
	} else {

		// Kick off part getter
		partChan := make(chan []products.Part)
		go l.LoadParts(ctx, partChan, page, count)

		if l.Vehicle.Submodel == "" { // Get Submodels
			if err := l.GetSubmodels(ctx); err != nil {
				return nil, errors.New("Trouble getting submodels for vehicle lookup")
			}
		} else { // Get configurations
			if err := l.GetConfigurations(ctx); err != nil {
				return nil, errors.New("Trouble getting configurations for vehicle lookup")
			}
		}

		select {
		case parts := <-partChan:
			if len(parts) > 0 {
				l.Parts = parts
				l.Filter, _ = apifilter.PartFilter(ctx, l.Parts, nil)
			}
		case <-time.After(5 * time.Second):

		}
	}

	return l, nil
}

// Parses the vehicle data out of the request
// body. It will first check for Content-Type as
// JSON and parse accordingly.
func LoadVehicle(r *http.Request) (v products.Vehicle) {
	defer r.Body.Close()

	if strings.Contains(strings.ToLower(r.Header.Get("Content-Type")), "json") {
		if data, err := ioutil.ReadAll(r.Body); err == nil || len(data) > 0 {
			err = json.Unmarshal(data, &v)
			if err == nil && v.Base.Year > 0 {
				return
			}
		}
	}

	// Get vehicle year
	y_str := r.FormValue("year")
	if y_str == "" {
		return
	}
	v.Base.Year, _ = strconv.Atoi(y_str)
	if v.Base.Year == 0 {
		return
	}
	delete(r.Form, "year")

	// Get vehicle make
	v.Base.Make = r.FormValue("make")
	if v.Base.Make == "" {
		return
	}
	delete(r.Form, "make")

	// Get vehicle model
	v.Base.Model = r.FormValue("model")
	if v.Base.Model == "" {
		return
	}
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Submodel = r.FormValue("submodel")
	if v.Submodel == "" {
		return
	}
	delete(r.Form, "submodel")
	delete(r.Form, "page")
	delete(r.Form, "count")

	// Get vehicle configuration options
	for key, opt := range r.Form {
		ignore := false
		for _, param := range ignoredFormParams {
			if param == strings.ToLower(key) {
				ignore = true
				break
			}
		}
		if !ignore && len(opt) > 0 {
			conf := products.Configuration{
				Key:   key,
				Value: opt[0],
			}
			v.Configurations = append(v.Configurations, conf)
		}
	}

	return
}

func GetVehicle(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	baseId, err := strconv.Atoi(r.FormValue("base"))
	if err != nil {
		return nil, err
	}

	subId, err := strconv.Atoi(r.FormValue("sub"))
	if err != nil {
		return nil, err
	}
	configVals := r.FormValue("configs")
	configs := strings.Split(configVals, ",")

	return vehicle.GetVehicle(ctx, baseId, subId, configs)
}

func Inquire(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	defer r.Body.Close()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil || len(data) == 0 {
		return nil, errors.New("missed payload")
	}

	var i products.VehicleInquiry
	err = json.Unmarshal(data, &i)
	if err != nil {
		return nil, errors.New("bad payload")
	}

	err = i.Push(ctx)
	if err != nil {
		return nil, errors.New("failed submission")
	}

	i.SendEmail(ctx)

	return true, nil
}
