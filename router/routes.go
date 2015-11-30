package router

import (
	"github.com/curt-labs/API/controllers/category"
	"net/http"

	"github.com/curt-labs/API/controllers"
	// "github.com/curt-labs/API/controllers/apiKeyType"
	// "github.com/curt-labs/API/controllers/applicationGuide"
	// "github.com/curt-labs/API/controllers/blog"
	// "github.com/curt-labs/API/controllers/brand"
	// "github.com/curt-labs/API/controllers/cache"
	// "github.com/curt-labs/API/controllers/cart"
	// "github.com/curt-labs/API/controllers/cartIntegration"
	// "github.com/curt-labs/API/controllers/category"
	// "github.com/curt-labs/API/controllers/contact"
	// "github.com/curt-labs/API/controllers/customer"
	// "github.com/curt-labs/API/controllers/dealers"
	// "github.com/curt-labs/API/controllers/faq"
	// "github.com/curt-labs/API/controllers/forum"
	// "github.com/curt-labs/API/controllers/geography"
	// "github.com/curt-labs/API/controllers/landingPages"
	// "github.com/curt-labs/API/controllers/middleware"
	// "github.com/curt-labs/API/controllers/news"
	// "github.com/curt-labs/API/controllers/part"
	// "github.com/curt-labs/API/controllers/salesrep"
	// "github.com/curt-labs/API/controllers/search"
	// "github.com/curt-labs/API/controllers/showcase"
	// "github.com/curt-labs/API/controllers/site"
	// "github.com/curt-labs/API/controllers/techSupport"
	// "github.com/curt-labs/API/controllers/testimonials"
	// "github.com/curt-labs/API/controllers/vehicle"
	// "github.com/curt-labs/API/controllers/videos"
	// "github.com/curt-labs/API/controllers/vinLookup"
	// "github.com/curt-labs/API/controllers/warranty"
	// "github.com/curt-labs/API/controllers/webProperty"
)

const (
	// PUBLIC_ENDPOINT Anyone can access
	PUBLIC_ENDPOINT = "PUBLIC"

	// SHOPIFY_ACCOUNT_ENDPOINT Binds with authentication based off token
	SHOPIFY_ACCOUNT_ENDPOINT = "SHOPIFY_ACCOUNT"

	// SHOPIFY_ACCOUNT_ENDPOINT Binds without authentication
	SHOPIFY_ACCOUNT_LOGIN_ENDPOINT = "SHOPIFY_ACCOUNT_LOGIN"

	// SHOPIFY_ACCOUNT_ENDPOINT Binds with full account authentication based off token
	SHOPIFY_ENDPOINT = "SHOPIFY"

	// KEYED_ENDPOINT Typical CURT authentication
	KEYED_ENDPOINT = "KEYED"
)

// Route ...
type Route struct {
	Name       string
	Method     string
	Pattern    string
	Middleware string
	Handler    http.HandlerFunc
}

var routes = []Route{
	Route{
		"Index",
		"GET",
		"/",
		PUBLIC_ENDPOINT,
		controllers.Index,
	},
	Route{
		"Status Checker",
		"GET",
		"/status",
		PUBLIC_ENDPOINT,
		controllers.Status,
	},
	Route{
		"Get API Key Typs",
		"GET",
		"/apiKeyTypes",
		PUBLIC_ENDPOINT,
		controllers.Status,
	},
	Route{
		"Get Category Tree",
		"GET",
		"/category",
		KEYED_ENDPOINT,
		category_ctlr.GetCategoryTree,
	},
	Route{
		"Get Category",
		"GET",
		"/category/:id",
		KEYED_ENDPOINT,
		category_ctlr.GetCategory,
	},
	Route{
		"Get Category Parts",
		"GET",
		"/category/:id/parts",
		KEYED_ENDPOINT,
		category_ctlr.GetCategoryParts,
	},
}
