package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"

	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/customer"
)

const (
	apiContext = "API_CONTEXT"
	respObject = "RESPONSE_OBJECT"
)

// APIContext Holds all the possible globals that we are going to want
// to use throughout the request lifecycle.
type APIContext struct {
	DB          *sql.DB
	Session     *mgo.Session
	Encoder     interface{}
	Params      httprouter.Params
	DataContext *apicontext.DataContext
}

// APIHandler Will delegate requests off the defined middleware and finally
// to the appropriate request endpoint.
type APIHandler struct {

	// BeforeFuncs A series a middleware that gets executed before
	// endpoint handlers
	BeforeFuncs []http.Handler

	// AfterFuncs A series a middleware that gets executed after
	// endpoint handlers
	AfterFuncs []func(http.Handler) http.Handler

	// H Defines a function definition for Object-Oriented handlers
	H func(*APIContext, http.ResponseWriter, *http.Request) (interface{}, error)

	// S Defines a function definition for a static endpoint, great
	// for uptime checks, redirects, direct ouput, etc. (Bypasses all middleware)
	S func(*APIContext, http.ResponseWriter, *http.Request)
}

func (fn APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if fn.H == nil && fn.S == nil {
		apierror.GenerateError("There hasn't been a handler declared for this route", nil, w, r, http.StatusInternalServerError)
		return
	}

	ctx := &APIContext{
		Params: ps,
	}

	if fn.S != nil {
		fn.S(ctx, w, r)
		return
	}

	context.Set(r, apiContext, ctx)

	for _, bf := range fn.BeforeFuncs {
		bf.ServeHTTP(w, r)
	}

	ctx = context.Get(r, apiContext).(*APIContext)

	obj, _ := fn.H(ctx, w, r)

	// context.Clear(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)

	return
}

// Wrap Wraps APIHandler into httprouter.Handle
func Wrap(h APIHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r, ps)
		// Gzip(h)
	})
}

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
