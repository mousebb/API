package router

import (
	"github.com/curt-labs/API/controllers"
	"github.com/curt-labs/API/controllers/apiKeyType"
	"github.com/curt-labs/API/controllers/applicationGuide"
	"github.com/curt-labs/API/controllers/brand"
	"github.com/curt-labs/API/controllers/cache"
	"github.com/curt-labs/API/controllers/cartIntegration"
	"github.com/curt-labs/API/controllers/category"
	"github.com/curt-labs/API/controllers/customer"
	"github.com/curt-labs/API/controllers/part"
	"github.com/curt-labs/API/controllers/search"
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
	Name    string
	Method  string
	Pattern string
	Handler middleware.APIHandler
}

var noAuth = []middleware.Middleware{
	middleware.WrapMiddleware(middleware.DB{}),
	middleware.WrapMiddleware(middleware.Logger{}),
}

var common = []middleware.Middleware{
	middleware.WrapMiddleware(middleware.DB{}),
	middleware.WrapMiddleware(middleware.Keyed{}),
	middleware.WrapMiddleware(middleware.Logger{}),
}

var commonPrivate = []middleware.Middleware{
	middleware.WrapMiddleware(middleware.DB{}),
	middleware.WrapMiddleware(middleware.Keyed{
		Type: "PRIVATE",
	}),
	middleware.WrapMiddleware(middleware.Logger{}),
}

var commonSudo = []middleware.Middleware{
	middleware.WrapMiddleware(middleware.DB{}),
	middleware.WrapMiddleware(middleware.Keyed{
		Type: "PRIVATE",
		Sudo: true,
	}),
	middleware.WrapMiddleware(middleware.Logger{}),
}

var routes = []Route{

	// Static handlers
	Route{"Index", "GET", "/", middleware.APIHandler{S: controllers.Index}},
	Route{"Status Checker", "GET", "/status", middleware.APIHandler{S: controllers.Status}},

	// API Key Management
	Route{"Get API Key Types", "GET", "/api/keys/types", middleware.APIHandler{H: apiKeyType.GetAPIKeyTypes, Middleware: common}},

	// Application Guides
	Route{"Get Application Guides by WebSite", "GET", "/applicationGuide/:id/website", middleware.APIHandler{H: applicationGuide.GetApplicationGuidesByWebsite, Middleware: common}},
	Route{"Get Application Guide", "GET", "/applicationGuide/:id", middleware.APIHandler{H: applicationGuide.GetApplicationGuide, Middleware: common}},

	// Category Endpoints
	Route{"Get Category Tree", "GET", "/category", middleware.APIHandler{H: categoryCtlr.GetCategoryTree, Middleware: common}},
	Route{"Get Category", "GET", "/category/:id", middleware.APIHandler{H: categoryCtlr.GetCategory, Middleware: common}},
	Route{"Get Category Parts", "GET", "/category/:id/parts", middleware.APIHandler{H: categoryCtlr.GetCategoryParts, Middleware: common}},

	// Brand information
	Route{"Get Brands", "GET", "/brands", middleware.APIHandler{H: brandCtlr.GetAllBrands, Middleware: common}},
	Route{"Get Brand", "GET", "/brands/:id", middleware.APIHandler{H: brandCtlr.GetBrand, Middleware: common}},

	// Customer Pricing
	Route{"Get Pricing", "GET", "/pricing", middleware.APIHandler{H: pricingCtlr.GetPricing, Middleware: commonPrivate}},
	Route{"Get Pricing Count", "GET", "/pricing/count", middleware.APIHandler{H: pricingCtlr.GetPricingCount, Middleware: commonPrivate}},
	Route{"Get All Part Prices", "GET", "/pricing/part", middleware.APIHandler{H: pricingCtlr.GetAllPartPrices, Middleware: commonPrivate}},
	Route{"Insert Part Price", "POST", "/pricing/part", middleware.APIHandler{H: pricingCtlr.CreatePrice, Middleware: commonPrivate}},
	Route{"Update Part Price", "PUT", "/pricing/part", middleware.APIHandler{H: pricingCtlr.UpdatePrice, Middleware: commonPrivate}},
	Route{"Get Prices by Part ID", "GET", "/pricing/part/:part", middleware.APIHandler{H: pricingCtlr.GetPartPricesByPartID, Middleware: commonPrivate}},
	Route{"Get All Price Types", "GET", "/pricing/types", middleware.APIHandler{H: pricingCtlr.GetAllPriceTypes, Middleware: commonPrivate}},
	Route{"Reset All Prices to Map", "POST", "/pricing/reset/map", middleware.APIHandler{H: pricingCtlr.ResetAllToMap, Middleware: commonPrivate}},
	Route{"Global Percentage Based Reset", "POST", "/pricing/global/:type/:percentage", middleware.APIHandler{H: pricingCtlr.Global, Middleware: commonPrivate}},
	Route{"Upload Pricing", "POST", "/upload", middleware.APIHandler{H: pricingCtlr.Upload, Middleware: commonPrivate}},
	Route{"Download Pricing", "POST", "/download", middleware.APIHandler{S: pricingCtlr.Download, Middleware: commonPrivate}},

	// Customer Management
	Route{"Get Customer", "GET", "/customer", middleware.APIHandler{H: customerCtlr.GetCustomer, Middleware: commonPrivate}},
	Route{"Get User", "GET", "/customer/user/key/:key", middleware.APIHandler{H: customerCtlr.GetUser, Middleware: commonSudo}},
	Route{"Get User By Identifier", "GET", "/customer/user/id/:id", middleware.APIHandler{H: customerCtlr.GetUserByIdentifier, Middleware: commonSudo}},
	Route{"Update User", "POST", "/customer/user", middleware.APIHandler{H: customerCtlr.AddUser, Middleware: commonSudo}},
	Route{"Update User", "PUT", "/customer/user", middleware.APIHandler{H: customerCtlr.UpdateUser, Middleware: commonPrivate}},
	Route{"Update User By Identifier", "PUT", "/customer/user/:id", middleware.APIHandler{H: customerCtlr.UpdateUser, Middleware: commonPrivate}},
	Route{"Authenticate User", "POST", "/customer/user/auth", middleware.APIHandler{H: customerCtlr.Authenticate, Middleware: noAuth}},
	Route{"Get User By Key", "GET", "/customer/user", middleware.APIHandler{H: customerCtlr.GetUserByKey, Middleware: commonPrivate}},

	// Cache Management
	Route{"Get Cache Keys", "GET", "/cache/keys", middleware.APIHandler{H: cache.GetKeys, Middleware: common}},
	Route{"Get Cache By Key", "GET", "/cache/key", middleware.APIHandler{H: cache.GetByKey, Middleware: common}},
	Route{"Delete By Key", "DELETE", "/cache/keys", middleware.APIHandler{H: cache.DeleteKey, Middleware: common}},

	// Product Management
	Route{"Get Featured Parts", "GET", "/parts/featured", middleware.APIHandler{H: partCtlr.Featured, Middleware: common}},
	Route{"Get Latest Parts", "GET", "/parts/latest", middleware.APIHandler{H: partCtlr.Latest, Middleware: common}},
	Route{"Get All Identifiers", "GET", "/parts/identifiers", middleware.APIHandler{H: partCtlr.Identifiers, Middleware: common}},
	Route{"Get Part Vehicles", "GET", "/part/:part/vehicles", middleware.APIHandler{H: partCtlr.Vehicles, Middleware: common}},
	Route{"Get Part", "GET", "/part/:part", middleware.APIHandler{H: partCtlr.Get, Middleware: common}},
	Route{"Get Part Attributes", "GET", "/part/:part/attributes", middleware.APIHandler{H: partCtlr.Attributes, Middleware: common}},
	Route{"Get Part Reviews", "GET", "/part/:part/reviews", middleware.APIHandler{H: partCtlr.ActiveApprovedReviews, Middleware: common}},
	Route{"Get Part Categories", "GET", "/part/:part/categories", middleware.APIHandler{H: partCtlr.Categories, Middleware: common}},
	Route{"Get Part Content", "GET", "/part/:part/content", middleware.APIHandler{H: partCtlr.GetContent, Middleware: common}},
	Route{"Get Part Images", "GET", "/part/:part/images", middleware.APIHandler{H: partCtlr.Images, Middleware: common}},
	// Route{"Get Part Installation Sheet", "GET", "/part/:part.pdf", middleware.APIHandler{S: partCtlr.InstallSheet, Middleware: common}},
	Route{"Get Part Packages", "GET", "/part/:part/packages", middleware.APIHandler{H: partCtlr.Packaging, Middleware: common}},
	Route{"Get Part Pricing", "GET", "/part/:part/pricing", middleware.APIHandler{H: partCtlr.Prices, Middleware: common}},
	Route{"Get Related Parts", "GET", "/part/:part/related", middleware.APIHandler{H: partCtlr.GetRelated, Middleware: common}},
	Route{"Get Videos", "GET", "/part/:part/videos", middleware.APIHandler{H: partCtlr.Videos, Middleware: common}},
	// Route{"Get Part with Vehicle", "GET", "/part/:part/vehicle/:year/:make/:model/:submodel", middleware.APIHandler{H: partCtlr.GetWithVehicle, Middleware: common}},
	Route{"Get Part with Vehicle Config", "GET", "/part/:part/vehicle/:year/:make/:model/:submodel/:config", middleware.APIHandler{H: partCtlr.GetWithVehicle, Middleware: common}},

	Route{"Get All Parts", "GET", "/parts", middleware.APIHandler{H: partCtlr.All, Middleware: common}},

	// Search
	Route{"Search", "GET", "/search/:term", middleware.APIHandler{H: searchCtlr.Search, Middleware: common}},
}
