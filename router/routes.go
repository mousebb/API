package router

import (
	"github.com/curt-labs/API/controllers"
	"github.com/curt-labs/API/controllers/apiKeyType"
	"github.com/curt-labs/API/controllers/applicationGuide"
	"github.com/curt-labs/API/controllers/brand"
	"github.com/curt-labs/API/controllers/cache"
	"github.com/curt-labs/API/controllers/cartIntegration"
	"github.com/curt-labs/API/controllers/category"
	"github.com/curt-labs/API/controllers/part"
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
	middleware.WrapMiddleware(middleware.DB{}),
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

	// Product Management
	Route{"Get Featured Parts", "GET", "/parts/featured", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Featured, Middleware: common}},
	Route{"Get Latest Parts", "GET", "/parts/latest", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Latest, Middleware: common}},
	Route{"Get All Identifiers", "GET", "/parts/identifiers", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Identifiers, Middleware: common}},
	Route{"Get Part Vehicles", "GET", "/parts/:part/vehicles", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Vehicles, Middleware: common}},
	Route{"Get Part Attributes", "GET", "/parts/:part/attributes", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Attributes, Middleware: common}},
	Route{"Get Part Reviews", "GET", "/parts/:part/reviews", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.ActiveApprovedReviews, Middleware: common}},
	Route{"Get Part Categories", "GET", "/parts/:part/categories", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Categories, Middleware: common}},
	Route{"Get Part Content", "GET", "/parts/:part/content", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.GetContent, Middleware: common}},
	Route{"Get Part Images", "GET", "/parts/:part/images", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Images, Middleware: common}},
	Route{"Get Part Installation Sheet", "GET", "/parts/:part.pdf", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.InstallSheet, Middleware: common}},
	Route{"Get Part Packages", "GET", "/parts/:part/packages", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Packaging, Middleware: common}},
	Route{"Get Part Pricing", "GET", "/parts/:part/pricing", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Prices, Middleware: common}},
	Route{"Get Related Parts", "GET", "/parts/:part/related", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.GetRelated, Middleware: common}},
	Route{"Get Videos", "GET", "/parts/:part/videos", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.Videos, Middleware: common}},
	Route{"Get Part with Vehicle", "GET", "/parts/:part/:year/:make/:model/:submodel", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.GetWithVehicle, Middleware: common}},
	Route{"Get Part with Vehicle Config", "GET", "/parts/:part/:year/:make/:model/:submodel/:config", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.GetWithVehicle, Middleware: common}},
	Route{"Get Part", "GET", "/parts/:part", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.PartNumber, Middleware: common}},
	Route{"Get All Parts", "GET", "/parts", KEYED_ENDPOINT, middleware.APIHandler{H: partCtlr.All, Middleware: common}},
}
