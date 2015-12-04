package router

import (
	"reflect"

	"github.com/curt-labs/API/middleware"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
}

func NewRouter() *Router {
	return &Router{httprouter.New()}
}

func (r *Router) HandleRoute(method string, pattern string, handler middleware.APIHandler) {
	r.Router.Handle(method, pattern, middleware.Wrap(handler))
}

// New ...
func New() *Router {
	r := NewRouter()
	r.Router.RedirectTrailingSlash = true

	for _, route := range routes {
		for _, m := range route.Handler.Middleware {
			nm := reflect.TypeOf(m).Name()
			if nm == "Keyed" {
				m.(middleware.Keyed).Type = "PRIVATE"
			}
		}
		r.HandleRoute(route.Method, route.Pattern, route.Handler)
		// switch route.Middleware {
		// case PUBLIC_ENDPOINT:
		// 	r.HandleRoute(route.Method, route.Pattern, route.Handler)
		// case SHOPIFY_ACCOUNT_LOGIN_ENDPOINT:
		// 	r.HandleRoute(route.Method, route.Pattern, route.Handler)
		// // case SHOPIFY_ACCOUNT_ENDPOINT:
		// // 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.ShopifyAccount, route.Handler))
		// // case SHOPIFY_ENDPOINT:
		// // 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.Shopify, route.Handler))
		// case KEYED_ENDPOINT:
		// 	r.HandleRoute(route.Method, route.Pattern, route.Handler)
		// }

	}

	return r
}
