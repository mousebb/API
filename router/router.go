package router

import (
	"github.com/curt-labs/API/middleware"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
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
			router.Handler(route.Method, route.Pattern, middleware.Wrapper(route.Name, route.Method, route.Handler))
			// router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, route.Handler))
			// case SHOPIFY_ACCOUNT_LOGIN_ENDPOINT:
			// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.ShopifyAccountLogin, route.Handler))
			// case SHOPIFY_ACCOUNT_ENDPOINT:
			// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.ShopifyAccount, route.Handler))
			// case SHOPIFY_ENDPOINT:
			// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.Shopify, route.Handler))
			// case KEYED_ENDPOINT:
			// 	router.Handler(route.Method, route.Pattern, middleware.Chain(route.Name, route.Method, middleware.Wrapper, middleware.Keyed, route.Handler))
		}

	}

	return router
}
