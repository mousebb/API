package vehicle

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"

	"net/http"
)

func CurtLookup(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var v products.CurtVehicle

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

	cl := products.CurtLookup{
		CurtVehicle: v,
	}

	var err error
	if v.Year == "" {
		err = cl.GetYears(ctx)
	} else if v.Make == "" {
		err = cl.GetMakes(ctx)
	} else if v.Model == "" {
		err = cl.GetModels(ctx)
	} else {
		err = cl.GetStyles(ctx)
		if err != nil {
			return nil, err
		}
		err = cl.GetParts(ctx)
	}

	return cl, err
}
