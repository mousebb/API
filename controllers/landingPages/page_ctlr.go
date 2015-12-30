package landingPage

import (
	//"encoding/json"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/landingPages"
)

func Get(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var lp landingPage.LandingPage
	var err error
	id, err := strconv.Atoi(ctx.Params.ByName("id"))

	if err != nil {
		return nil, err
	}

	lp.Id = id
	err = lp.Get(ctx)
	return lp, err
}
