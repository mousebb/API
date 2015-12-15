package applicationGuide

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/applicationGuide"
)

// GetApplicationGuide ...
func GetApplicationGuide(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var ag applicationGuide.ApplicationGuide

	ag.ID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return nil, err
	}

	err = ag.Get(ctx)
	if err != nil {
		return nil, err
	}

	return ag, nil
}

// GetApplicationGuidesByWebsite ...
func GetApplicationGuidesByWebsite(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var ag applicationGuide.ApplicationGuide
	var err error
	ag.Website.ID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return nil, fmt.Errorf("%s", "failed to parse website identifier")
	}

	return ag.GetBySite(ctx)
}
