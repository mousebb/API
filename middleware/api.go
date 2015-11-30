package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
	"github.com/gorilla/context"
)

type ApiHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (fn ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	context.Set(r, "params", ps)

	obj, err := fn(w, r)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(obj)
	return
}

// Keyed ...
func Keyed(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var err error

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
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}
		// go user.LogApiRequest(r)

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
		// var dtx apicontext.DataContext
		dtx := &apicontext.DataContext{
			APIKey:     apiKey,
			BrandID:    brandID,
			WebsiteID:  websiteID,
			UserID:     user.Id, //current authenticated user
			CustomerID: user.CustomerID,
			Globals:    nil,
		}
		err = dtx.GetBrandsArrayAndString(apiKey, brandID)
		if err != nil {
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		context.Set(r, "request_ctx", dtx)

		h.ServeHTTP(rw, r)
	})
}
