// Package vehicle Allows request operations for a vehicle lookup.
// @SubApi Vehicle Application Lookup [/vehicle]
package vehicle

import (
	"encoding/json"
	"net/http"

	"github.com/curt-labs/API/helpers/rest"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
)

var (
	ignoredFormParams = []string{"key"}
)

// Query Takes in basic vehicle data and returns further qualifiers or an
// array of products.
// @Title Query
// @Description Query for Vehicle Application data
// @Accept application/x-www-form-urlencoded
func Query(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var va products.VehicleApplication
	var err error
	if rest.IsJsonRequest(r) {
		err = json.NewDecoder(r.Body).Decode(&va)
		if err != nil {
			return nil, err
		}
	} else {
		va.Year = r.FormValue("year")
		va.Make = r.FormValue("make")
		va.Model = r.FormValue("model")
	}

	va, err = products.Query(
		ctx,
		va.Year,
		va.Make,
		va.Model,
		r.URL.Query().Get("category"),
	)
	if err != nil {
		return nil, err
	}

	return va, nil
}
