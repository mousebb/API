package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
)

type Handler func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error

// Wrapper Wraps all endpoints in the generic middleware
func Wrapper(name string, method string, handlers ...Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		for _, handler := range handlers {
			err := handler(w, r, ps)
			if err != nil {
				apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
				return
			}
		}
	}
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
