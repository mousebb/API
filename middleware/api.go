package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"github.com/curt-labs/API/helpers/error"
)

const (
	apiContext = "API_CONTEXT"
	respObject = "RESPONSE_OBJECT"
)

// APIContext Holds all the possible globals that we are going to want
// to use throughout the request lifecycle.
type APIContext struct {
	Brand              int
	DB                 *sql.DB
	Session            *mgo.Session
	AriesSession       *mgo.Session
	AriesMongoDatabase string
	Encoder            interface{}
	Params             httprouter.Params
	DataContext        *DataContext
}

// Middleware Gives the ability to add Middleware capability to APIHandler
// that supports before and deferred after functionality.
type Middleware struct {
	H     http.Handler
	Defer bool
}

// APIHandler Will delegate requests off the defined middleware and finally
// to the appropriate request endpoint.
type APIHandler struct {

	// BeforeFuncs A series a middleware that gets executed before
	// endpoint handlers
	Middleware []Middleware

	// AfterFuncs A series a middleware that gets executed after
	// endpoint handlers
	AfterFuncs []func(http.Handler) http.Handler

	// H Defines a function definition for Object-Oriented handlers
	H func(*APIContext, http.ResponseWriter, *http.Request) (interface{}, error)

	// S Defines a function definition for a static endpoint, great
	// for uptime checks, redirects, direct ouput, etc. (Bypasses all middleware)
	S func(*APIContext, http.ResponseWriter, *http.Request)
}

// ServeHTTP For interfacing http.HandlerFunc
func (fn APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if fn.H == nil && fn.S == nil {
		apierror.GenerateError("There hasn't been a handler declared for this route", nil, w, r, http.StatusInternalServerError)
		return
	}

	ctx := &APIContext{
		Params: ps,
	}

	context.Set(r, apiContext, ctx)

	for _, m := range fn.Middleware {
		if m.H == nil {
			continue
		}
		if m.Defer {
			defer m.H.ServeHTTP(w, r)
		} else {
			rec := httptest.NewRecorder()
			m.H.ServeHTTP(rec, r)
			if rec.Code != 200 {
				err := fmt.Errorf("%s", rec.Body.String())
				w.WriteHeader(rec.Code)
				w.Write([]byte(err.Error()))

				return
			}

		}
	}

	ctx = context.Get(r, apiContext).(*APIContext)

	if fn.S != nil {
		fn.S(ctx, w, r)
		return
	}

	obj, err := fn.H(ctx, w, r)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	context.Clear(r)

	err = Encode(r, w, obj)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
		return
	}

	return
}

// Wrap Wraps APIHandler into httprouter.Handle
func Wrap(h APIHandler) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		h.ServeHTTP(w, r, ps)
	})
}

// WrapMiddleware Convenience function around WrapDeferredMiddleware
// defaulting the deferred option to false.
func WrapMiddleware(h http.Handler) Middleware {
	return WrapDeferredMiddleware(h, false)
}

// WrapDeferredMiddleware Converts http.HandlerFunc to Middleware with defered
// designation.
func WrapDeferredMiddleware(h http.Handler, def bool) Middleware {
	return Middleware{h, def}
}
