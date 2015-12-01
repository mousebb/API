package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
	"github.com/gorilla/context"
)

type ApiContext struct {
	db      *sql.DB
	session *mgo.Session
	encoder interface{}
}

type ApiHandler struct {
	*ApiContext
	H func(*ApiContext, http.ResponseWriter, *http.Request) (interface{}, error)
	S func(*ApiContext, http.ResponseWriter, *http.Request)
}

func (fn ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if fn.H == nil && fn.S == nil {
		apierror.GenerateError("There hasn't been a handler declared for this route", nil, w, r, http.StatusInternalServerError)
		return
	}
	context.Set(r, "params", ps)

	if fn.S != nil {
		fn.S(fn.ApiContext, w, r)
		return
	}

	obj, err := fn.H(fn.ApiContext, w, r)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(obj)
}

func Wrap(h ApiHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r, ps)
	})
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
