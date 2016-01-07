package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/error"
	"github.com/gorilla/context"
)

// Keyed http.HandlerFunc middleware that will validate
// API key usage
type Keyed struct {
	Type    string
	Sudo    bool
	Error   error
	Status  int
	handler http.Handler
}

func (kh Keyed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx := context.Get(r, apiContext).(*APIContext)
	if ctx == nil {
		ctx = &APIContext{}
	}
	if ctx.DataContext == nil {
		ctx.DataContext = &DataContext{}
	}

	qs := r.URL.Query()
	apiKey := qs.Get("key")
	brand := qs.Get("brandID")
	website := qs.Get("websiteID")

	//handles api key
	if apiKey == "" {
		apiKey = r.FormValue("key")
	}
	if apiKey == "" {
		apiKey = r.Header.Get("key")
	}
	if apiKey == "" {
		err = fmt.Errorf("%s", "No API Key Supplied.")
	}

	//handles branding
	if brand == "" {
		brand = r.FormValue("brandID")
	}
	if brand == "" {
		brand = r.Header.Get("brandID")
	}
	ctx.DataContext.BrandID, _ = strconv.Atoi(brand)

	//handles websiteID
	if website == "" {
		website = r.FormValue("websiteID")
	}
	if website == "" {
		website = r.Header.Get("websiteID")
	}
	ctx.DataContext.WebsiteID, _ = strconv.Atoi(website)

	//load brands in dtx
	//returns our data context...shared amongst controllers
	err = ctx.BuildDataContext(apiKey, kh.Type, kh.Sudo)
	if err != nil {
		err = fmt.Errorf("failed to authenticate for %s key %s", kh.Type, apiKey)
		apierror.GenerateError(err.Error(), err, w, r, http.StatusUnauthorized)
		return
	}

	context.Set(r, apiContext, ctx)

	return
}
