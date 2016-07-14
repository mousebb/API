package partCtlr

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/rest"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
)

// Identifiers Returns a slice of distinct part numbers.
func Identifiers(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return products.Identifiers(ctx)
}

// All Returns a slice of all Part.
func All(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	page := 0
	count := 10
	qs := r.URL.Query()

	if qs.Get("page") != "" {
		if pg, err := strconv.Atoi(qs.Get("page")); err == nil {
			if pg == 0 {
				pg = 1
			}
			page = pg - 1
		}
	}
	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 500 {
				return nil, fmt.Errorf("maximum request size is 500, you requested: %d", ct)
			}
			count = ct
		}
	}

	return products.All(ctx, page, count)
}

// Featured Returns a given amount of featured Part.
func Featured(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				return nil, fmt.Errorf("maximum request size is 50, you requested: %d", ct)
			}
			count = ct
		}
	}

	return products.Featured(ctx, count)
}

// Latest Returns the latest slice of Part.
func Latest(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	count := 10
	qs := r.URL.Query()

	if qs.Get("count") != "" {
		if ct, err := strconv.Atoi(qs.Get("count")); err == nil {
			if ct > 50 {
				return nil, fmt.Errorf("maximum request size is 50, you requested: %d", ct)
			}
			count = ct
		}
	}

	return products.Latest(ctx, count)
}

// Get retrieves a Part using the :part segment of the URL pattern.
// TODO: We should add logic for the year/make/model/style.
func Get(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		SKU: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx)
	return p, err
}

// GetAttributes retrieves an array of Attribute using the :part segment of the URL pattern.
func GetAttributes(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return products.GetAttributes(ctx, ctx.Params.ByName("part"))
}

// GetRelated Retrieves the related Part to a given Part.
func GetRelated(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	p := products.Part{
		SKU: ctx.Params.ByName("part"),
	}

	return p.GetRelated(ctx)
}

// GetVehicles Retrieves the related Part to a given Part.
func GetVehicles(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return products.GetVehicles(ctx, ctx.Params.ByName("part"))
}

// GetContent Retrieves the related Part to a given Part.
// func GetContent(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetContent(ctx, ctx.Params.ByName("part"))
// }
//
// // GetImages Retrieves the related Part to a given Part.
// func GetImages(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetImages(ctx, ctx.Params.ByName("part"))
// }
//
// // GetPackaging Retrieves the related Part to a given Part.
// func GetPackaging(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetPackaging(ctx, ctx.Params.ByName("part"))
// }
//
// // GetReviews Retrieves the related Part to a given Part.
// func GetReviews(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetReviews(ctx, ctx.Params.ByName("part"))
// }
//
// // GetVideos Retrieves the related Part to a given Part.
// func GetVideos(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetVideos(ctx, ctx.Params.ByName("part"))
// }
//
// // GetCategories Retrieves the related Part to a given Part.
// func GetCategories(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	return products.GetCategories(ctx, ctx.Params.ByName("part"))
// }
//
// // GetPrices Retrieves the related Part to a given Part.
// func GetPrices(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
// 	sku := ctx.Params.ByName("part")
// 	prices, err := products.GetPrices(ctx, sku)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	price, err := customer.GetCustomerPrice(ctx.DB, sku)
// 	if err == nil {
// 		custPrice := products.Price{
// 			Id:           0,
// 			PartId:       0,
// 			Type:         "Customer",
// 			Price:        price,
// 			Enforced:     false,
// 			DateModified: time.Now(),
// 		}
// 		prices = append(prices, custPrice)
// 	}
//
// 	return prices, nil
// }

// InstallSheet Takes the provided part number and outputs the assoicated installation
// sheet to the response.
func InstallSheet(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) {
	p := products.Part{
		SKU: ctx.Params.ByName("part"),
	}

	err := p.Get(ctx)
	if err != nil {
		apierror.GenerateError("failed to find product", err, rw, r)
		return
	}

	if p.InstallSheet == nil {
		apierror.GenerateError("no installation sheet for this part", err, rw, r, http.StatusNoContent)
		return
	}

	data, err := rest.GetPDF(p.InstallSheet.String())
	if err != nil {
		apierror.GenerateError("Error getting PDF", err, rw, r)
		return
	}

	rw.Header().Set("Content-Length", strconv.Itoa(len(data)))
	rw.Header().Set("Content-Type", "application/pdf")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rw.Header().Set("Access-Control-Allow-Headers", "Origin")
	rw.Write(data)
}
