package brandCtlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/brand"
)

// GetAllBrands Returns an a slice of Brand.
// /brands
func GetAllBrands(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return brand.GetAllBrands(ctx.DB)
}

// GetBrand Returns a specific Brand based off
// the :id paramter in the route.
// /brands/:id
func GetBrand(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	var br brand.Brand

	if br.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}
	err = br.Get(ctx.DB)

	return br, err
}
