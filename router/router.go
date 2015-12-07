package router

import (
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
		r.HandleRoute(route.Method, route.Pattern, route.Handler)
	}

	return r
}
