package router

import (
	"net/http"

	"github.com/curt-labs/API/middleware"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func (r *Router) HandleRoute(method string, pattern string, handler http.Handler) {
	r.Router.Handle(method, pattern, middleware.ApiHandler(handler))
}

// func wrapHandler(h http.Handler) httprouter.Handle {
// 	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 		obj, err := h(w, r, ps)
// 		if err != nil {
// 			apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
// 			return
// 		}

// 		json.NewEncoder(w).Encode(obj)

// 	}
// }

// New ...
func New() *Router {
	r := NewRouter()
	r.Router.RedirectTrailingSlash = true

	common := alice.New(context.ClearHandler)

	for _, route := range routes {
		switch route.Middleware {
		case PUBLIC_ENDPOINT:
			r.HandleRoute(route.Method, route.Pattern, common.ThenFunc(route.Handler))
		case SHOPIFY_ACCOUNT_LOGIN_ENDPOINT:
			r.HandleRoute(route.Method, route.Pattern, common.Append(middleware.ShopifyAccountLogin).ThenFunc(middleware.ApiHandler(route.Handler)))
		// case SHOPIFY_ACCOUNT_ENDPOINT:
		// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.ShopifyAccount, route.Handler))
		// case SHOPIFY_ENDPOINT:
		// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.Shopify, route.Handler))
		case KEYED_ENDPOINT:
			r.HandleRoute(route.Method, route.Pattern, common.Append(middleware.Keyed).ThenFunc(route.Handler))
		}

	}

	return r
}
