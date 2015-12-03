package router

import (
	"github.com/curt-labs/API/controllers"
	"github.com/curt-labs/API/controllers/apiKeyType"
	"github.com/curt-labs/API/controllers/applicationGuide"
	"github.com/curt-labs/API/controllers/brand"
	"github.com/curt-labs/API/controllers/cache"
	"github.com/curt-labs/API/controllers/cartIntegration"
	"github.com/curt-labs/API/controllers/category"
	"github.com/curt-labs/API/middleware"
)

const (
	// PUBLIC_ENDPOINT Anyone can access
	PUBLIC_ENDPOINT = "PUBLIC"

	// KEYED_ENDPOINT Typical CURT authentication
	KEYED_ENDPOINT = "KEYED"

	// KEYED_PRIVATE_ENDPOINT Requires a private API key
	KEYED_PRIVATE_ENDPOINT = "KEYED_PRIVATE"
)

// Route ...
type Route struct {
	Name       string
	Method     string
	Pattern    string
	Middleware string
	Handler    middleware.APIHandler
}

var common = []middleware.Middleware{
	middleware.WrapMiddleware(middleware.Mongo{}),
	middleware.WrapMiddleware(middleware.Keyed{}),
	middleware.WrapMiddleware(middleware.Logger{}),
}

var routes = []Route{

	// Static handlers
	Route{"Index", "GET", "/", PUBLIC_ENDPOINT, middleware.APIHandler{S: controllers.Index}},
	Route{"Status Checker", "GET", "/status", PUBLIC_ENDPOINT, middleware.APIHandler{S: controllers.Status}},

	// API Key Management
	Route{"Get API Key Types", "GET", "/api/keys/types", KEYED_ENDPOINT, middleware.APIHandler{H: apiKeyType.GetApiKeyTypes, Middleware: common}},

	// Application Guides
	Route{"Get Application Guides by WebSite", "GET", "/applicationGuide/:id/website", KEYED_ENDPOINT, middleware.APIHandler{H: applicationGuide.GetApplicationGuidesByWebsite, Middleware: common}},
	Route{"Get Application Guide", "GET", "/applicationGuide/:id", KEYED_ENDPOINT, middleware.APIHandler{H: applicationGuide.GetApplicationGuide, Middleware: common}},

	// Category Endpoints
	Route{"Get Category Tree", "GET", "/category", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategoryTree, Middleware: common}},
	Route{"Get Category", "GET", "/category/:id", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategory, Middleware: common}},
	Route{"Get Category Parts", "GET", "/category/:id/parts", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategoryParts, Middleware: common}},

	// Brand information
	Route{"Get Brands", "GET", "/brands", KEYED_ENDPOINT, middleware.APIHandler{H: brandCtlr.GetAllBrands, Middleware: common}},
	Route{"Get Brand", "GET", "/brands/:id", KEYED_ENDPOINT, middleware.APIHandler{H: brandCtlr.GetBrand, Middleware: common}},

	// Customer Pricing
	Route{"Get Pricing", "GET", "/pricing", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetPricing, Middleware: common}},
	Route{"Get Pricing Count", "GET", "/pricing/count", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetPricingCount, Middleware: common}},
	Route{"Get All Part Prices", "GET", "/pricing/part", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetAllPartPrices, Middleware: common}},
	Route{"Insert Part Price", "POST", "/pricing/part", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.CreatePrice, Middleware: common}},
	Route{"Update Part Price", "PUT", "/pricing/part", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.UpdatePrice, Middleware: common}},
	Route{"Get Prices by Part ID", "GET", "/pricing/part/:part", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetPartPricesByPartID, Middleware: common}},
	Route{"Get All Price Types", "GET", "/pricing/types", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetAllPriceTypes, Middleware: common}},
	Route{"Reset All Prices to Map", "POST", "/pricing/reset/map", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.ResetAllToMap, Middleware: common}},
	Route{"Global Percentage Based Reset", "POST", "/pricing/global/:type/:percentage", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.Global, Middleware: common}},
	Route{"Upload Pricing", "POST", "/upload", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.Upload, Middleware: common}},
	Route{"Download Pricing", "POST", "/download", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{S: pricingCtlr.Download, Middleware: common}},
	Route{"Get Pricing Paged", "GET", "/pricing/:page/:count", KEYED_PRIVATE_ENDPOINT, middleware.APIHandler{H: pricingCtlr.GetPricingPaged, Middleware: common}},

	// Cache Management
	Route{"Get Cache Keys", "GET", "/cache/keys", KEYED_ENDPOINT, middleware.APIHandler{H: cache.GetKeys, Middleware: common}},
	Route{"Get Cache By Key", "GET", "/cache/key", KEYED_ENDPOINT, middleware.APIHandler{H: cache.GetByKey, Middleware: common}},
	Route{"Delete By Key", "DELETE", "/cache/keys", KEYED_ENDPOINT, middleware.APIHandler{H: cache.DeleteKey, Middleware: common}},

	// Customer Management
}
