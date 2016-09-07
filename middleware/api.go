package middleware

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"github.com/curt-labs/API/helpers/beefwriter"
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
	Brand              int
	DB                 *sql.DB
	Session            *mgo.Session
	AriesSession       *mgo.Session
	AriesMongoDatabase string
	Encoder            interface{}
	Params             httprouter.Params
	DataContext        *customer.DataContext
	RequestStart       time.Time

	// Statuses are the product implementation statuses that this request
	// is looking to retrieve.
	Statuses []int
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
	var rw beefwriter.ResponseWriter
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		rw = beefwriter.NewResponseWriter(w)
	} else {
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(beefwriter.NewResponseWriter(w))
		defer gz.Close()
		rw = beefwriter.NewGzipResponseWriter(gz, w)
	}

	ctx := &APIContext{
		Params:       ps,
		RequestStart: time.Now(),
		Statuses:     []int{800, 900},
	}

	if r.Header.Get("X-Statuses") != "" {
		var statuses []int
		segs := strings.Split(r.Header.Get("X-Statuses"), ",")
		for _, seg := range segs {
			status, err := strconv.Atoi(seg)
			if err == nil {
				statuses = append(statuses, status)
			}
		}

		ctx.Statuses = statuses
	}

	context.Set(r, apiContext, ctx)

	for _, m := range fn.Middleware {
		if m.H == nil {
			continue
		}
		if m.Defer {
			defer m.H.ServeHTTP(rw, r)
		} else {
			rec := httptest.NewRecorder()
			m.H.ServeHTTP(rec, r)
			if rec.Code != 200 {
				err := fmt.Errorf("%s", rec.Body.String())
				rw.WriteHeader(rec.Code)
				rw.Write([]byte(err.Error()))

				return
			}

		}
	}

	ctx = context.Get(r, apiContext).(*APIContext)

	if fn.S != nil {
		fn.S(ctx, rw, r)
		return
	}

	obj, err := fn.H(ctx, rw, r)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r, http.StatusInternalServerError)
		return
	}

	context.Clear(r)

	err = Encode(r, rw, obj)
	if err != nil {
		apierror.GenerateError(err.Error(), err, rw, r, http.StatusInternalServerError)
		return
	}

	go Log(rw, *r, *ctx)

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
