package vehicle

import (
	"net/http"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
)

var (
	ignoredFormParams = []string{"key"}
)

func Query(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	va := products.VehicleApplication{
		Year:  r.FormValue("year"),
		Make:  r.FormValue("make"),
		Model: r.FormValue("model"),
	}

	va, err := products.Query(
		ctx,
		r.FormValue("year"),
		r.FormValue("make"),
		r.FormValue("model"),
	)
	if err != nil {
		return nil, err
	}

	return va, nil
}
