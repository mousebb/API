package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
	"github.com/gorilla/context"
)

// Keyed http.HandlerFunc middleware that will validate
// API key usage
type Keyed struct {
	http.Handler
}

func (kh Keyed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	ctx := context.Get(r, apiContext).(*APIContext)
	if ctx == nil {
		ctx = &APIContext{}
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

	//gets customer user from api key
	user, err := customer.GetCustomerUserFromKey(apiKey)
	if err != nil || user.Id == "" {
		err = fmt.Errorf("%s", "No User for this API Key.")
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	//handles branding
	var brandID int
	if brand == "" {
		brand = r.FormValue("brandID")
	}
	if brand == "" {
		brand = r.Header.Get("brandID")
	}
	brandID, _ = strconv.Atoi(brand)

	//handles websiteID
	if website == "" {
		website = r.FormValue("websiteID")
	}
	if website == "" {
		website = r.Header.Get("websiteID")
	}
	websiteID, _ := strconv.Atoi(website)

	//load brands in dtx
	//returns our data context...shared amongst controllers
	ctx.DataContext = &apicontext.DataContext{
		APIKey:     apiKey,
		BrandID:    brandID,
		WebsiteID:  websiteID,
		UserID:     user.Id, //current authenticated user
		CustomerID: user.CustomerID,
		Globals:    nil,
	}

	err = ctx.DataContext.GetBrandsArrayAndString(apiKey, brandID)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	context.Set(r, apiContext, ctx)
}
