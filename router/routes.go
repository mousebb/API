package router

import (
	"net/http"

	"github.com/curt-labs/API/middleware"
	"github.com/gorilla/context"

	"github.com/curt-labs/API/controllers"
	"github.com/curt-labs/API/controllers/apiKeyType"
	"github.com/curt-labs/API/controllers/applicationGuide"
	"github.com/curt-labs/API/controllers/category"
)

const (
	// PUBLIC_ENDPOINT Anyone can access
	PUBLIC_ENDPOINT = "PUBLIC"

	// SHOPIFY_ACCOUNT_ENDPOINT Binds with authentication based off token
	SHOPIFY_ACCOUNT_ENDPOINT = "SHOPIFY_ACCOUNT"

	// SHOPIFY_ACCOUNT_LOGIN_ENDPOINT Binds without authentication
	SHOPIFY_ACCOUNT_LOGIN_ENDPOINT = "SHOPIFY_ACCOUNT_LOGIN"

	// SHOPIFY_ENDPOINT Binds with full account authentication based off token
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
	Handler    middleware.APIHandler
}

var commonBefore = []http.Handler{middleware.Keyed{}}

// var commonBefore = []http.Handler{}
var commonAfter = []func(http.Handler) http.Handler{context.ClearHandler}

var routes = []Route{

	// Static handlers
	Route{"Index", "GET", "/", PUBLIC_ENDPOINT, middleware.APIHandler{S: controllers.Index}},
	Route{"Status Checker", "GET", "/status", PUBLIC_ENDPOINT, middleware.APIHandler{S: controllers.Status}},

	// API Key Management
	Route{"Get API Key Types", "GET", "/apiKeyTypes", KEYED_ENDPOINT, middleware.APIHandler{H: apiKeyType.GetApiKeyTypes, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},

	// Application Guides
	Route{"Get Application Guides by WebSite", "GET", "/applicationGuide/:id/website", KEYED_ENDPOINT, middleware.APIHandler{H: applicationGuide.GetApplicationGuidesByWebsite, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},
	Route{"Get Application Guide", "GET", "/applicationGuide/:id", KEYED_ENDPOINT, middleware.APIHandler{H: applicationGuide.GetApplicationGuide, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},

	// Category Endpoints
	Route{"Get Category Tree", "GET", "/category", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategoryTree, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},
	Route{"Get Category", "GET", "/category/:id", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategory, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},
	Route{"Get Category Parts", "GET", "/category/:id/parts", KEYED_ENDPOINT, middleware.APIHandler{H: categoryCtlr.GetCategoryParts, BeforeFuncs: commonBefore, AfterFuncs: commonAfter}},
}
