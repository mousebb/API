package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/cart"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2/bson"
)

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
