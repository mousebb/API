package router

import (
	"fmt"
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cart"
	"github.com/curt-labs/API/models/customer"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	rndr = render.New(render.Options{})
)

// New ...
func New() *httprouter.Router {
	router := httprouter.New()
	router.RedirectTrailingSlash = true

	for _, route := range routes {
		switch route.Middleware {
		case PUBLIC_ENDPOINT:
			router.Handler(route.Method, route.Pattern, ApiWrapper(route, route.Handler))
		case SHOPIFY_ACCOUNT_LOGIN_ENDPOINT:
			router.Handler(route.Method, route.Pattern, ApiWrapper(route, ShopifyAccountLogin(route.Handler)))
		case SHOPIFY_ACCOUNT_ENDPOINT:
			router.Handler(route.Method, route.Pattern, ApiWrapper(route, ShopifyAccount(route.Handler)))
		case SHOPIFY_ENDPOINT:
			router.Handler(route.Method, route.Pattern, ApiWrapper(route, Shopify(route.Handler)))
		case KEYED_ENDPOINT:
			router.Handler(route.Method, route.Pattern, ApiWrapper(route, Keyed(route.Handler)))
		}

	}

	return router
}

// ApiWrapper ...
func ApiWrapper(route Route, h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		cors.Default().Handler(h)
		context.Set(r, "renderer", rndr)

		h.ServeHTTP(rw, r)

		end := time.Now()
		fmt.Printf("%s\t%s_%s\t%s\n", route.Name, route.Method, r.URL.String(), end.Sub(start))
	})
}

// ShopifyAccount ...
func ShopifyAccount(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		shopID := r.URL.Query().Get("shop")
		var crt cart.Shop
		if bson.IsObjectIdHex(shopID) {
			crt.Id = bson.ObjectIdHex(shopID)
		}

		context.Set(r, "cart", &crt)

		h.ServeHTTP(rw, r)
	})
}

// ShopifyAccountLogin ...
func ShopifyAccountLogin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")
		token := strings.Replace(auth, "Bearer ", "", 1)
		var err error

		cust, err := cart.AuthenticateAccount(token)
		if err != nil {
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		shop := cart.Shop{
			Id: cust.ShopId,
		}

		if shop.Id.Hex() == "" {
			err = fmt.Errorf("error: %s", "invalid shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		if err := shop.Get(); err != nil {
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}
		if shop.Id.Hex() == "" {
			err = fmt.Errorf("error: %s", "invalid shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		context.Set(r, "shop", &shop)
		context.Set(r, "token", token)

		h.ServeHTTP(rw, r)
	})
}

// Shopify ...
func Shopify(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		qs := r.URL.Query()
		var shopID string
		var err error

		if qsID := qs.Get("shop"); qsID != "" {
			shopID = qsID
		} else if formID := r.FormValue("shop"); formID != "" {
			shopID = formID
		} else if headerID := r.Header.Get("shop"); headerID != "" {
			shopID = headerID
		}

		if shopID == "" {
			err = fmt.Errorf("error: %s", "you must provide a shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}
		if !bson.IsObjectIdHex(shopID) {
			err = fmt.Errorf("error: %s", "invalid shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}
		shop := cart.Shop{
			Id: bson.ObjectIdHex(shopID),
		}

		if shop.Id.Hex() == "" {
			err = fmt.Errorf("error: %s", "invalid shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		if err := shop.Get(); err != nil {
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}
		if shop.Id.Hex() == "" {
			err = fmt.Errorf("error: %s", "invalid shop identifier")
			apierror.GenerateError(err.Error(), err, rw, r)
			return
		}

		context.Set(r, "shop", &shop)

		h.ServeHTTP(rw, r)
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
