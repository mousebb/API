package geography

import (
	"net/http"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/geography"
)

func GetAllCountriesAndStates(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return geography.GetAllCountriesAndStates(ctx.DB)
}

func GetAllCountries(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return geography.GetAllCountries(ctx.DB)
}

func GetAllStates(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return geography.GetAllStates(ctx.DB)
}
