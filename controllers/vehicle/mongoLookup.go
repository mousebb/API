package vehicle

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"

	"net/http"
	"sort"
	"strconv"
)

func Collections(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return products.GetAriesVehicleCollections()
}

func Lookup(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var v products.NoSqlVehicle
	var collection string //e.g. interior/exterior

	//Get collection
	collection = r.FormValue("collection")
	delete(r.Form, "collection")

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	return products.FindVehicles(v, collection, dtx)
}

func ByCategory(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	collection := r.FormValue("collection")
	page, _ := strconv.Atoi(r.FormValue("page"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	var offset int
	if page == 0 {
		offset = 0
	} else if page == 1 {
		offset = 101
	} else {
		offset = page*limit + 1
	}

	return products.FindApplications(collection, offset, limit)
}

//Hack version that slowly traverses all the collection and aggregates results
func AllCollectionsLookup(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var v products.NoSqlVehicle

	//Get all collections
	cols, err := products.GetAriesVehicleCollections()
	if err != nil {
		return nil, err
	}

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	var collectionVehicleArray []products.NoSqlLookup

	for _, col := range cols {
		noSqlLookup, err := products.FindVehiclesWithParts(v, col, dtx)
		if err != nil {
			return nil, err
		}
		collectionVehicleArray = append(collectionVehicleArray, noSqlLookup)
	}
	l := makeLookupFrommanyLookups(collectionVehicleArray)

	return l, nil
}

//return parts for a vehicle(incl style) within a specific category
func AllCollectionsLookupCategory(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var v products.NoSqlVehicle
	noSqlLookup := make(map[string]products.NoSqlLookup)
	var err error

	// Get vehicle year
	v.Year = r.FormValue("year")
	delete(r.Form, "year")

	// Get vehicle make
	v.Make = r.FormValue("make")
	delete(r.Form, "make")

	// Get vehicle model
	v.Model = r.FormValue("model")
	delete(r.Form, "model")

	// // Get vehicle submodel
	v.Style = r.FormValue("style")
	delete(r.Form, "style")

	collection := r.FormValue("collection")
	if collection == "" {
		return products.FindVehiclesFromAllCategories(v, dtx)
	}

	return products.FindPartsFromOneCategory(v, collection, dtx)
}

func makeLookupFrommanyLookups(lookupArrays []products.NoSqlLookup) (l products.NoSqlLookup) {
	yearmap := make(map[string]string)
	makemap := make(map[string]string)
	modelmap := make(map[string]string)
	stylemap := make(map[string]string)
	partmap := make(map[int]products.Part)

	for _, lookup := range lookupArrays {
		for _, year := range lookup.Years {
			yearmap[year] = year
		}
		for _, mk := range lookup.Makes {
			makemap[mk] = mk
		}
		for _, model := range lookup.Models {
			modelmap[model] = model
		}
		for _, style := range lookup.Styles {
			stylemap[style] = style
		}
		for _, part := range lookup.Parts {
			partmap[part.ID] = part
		}
	}
	for year, _ := range yearmap {
		l.Years = append(l.Years, year)
	}
	for mk, _ := range makemap {
		l.Makes = append(l.Makes, mk)
	}
	for model, _ := range modelmap {
		l.Models = append(l.Models, model)
	}
	for style, _ := range stylemap {
		l.Styles = append(l.Styles, style)
	}
	for _, part := range partmap {
		l.Parts = append(l.Parts, part)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(l.Years)))
	sort.Strings(l.Makes)
	sort.Strings(l.Models)
	sort.Strings(l.Styles)

	return l
}
