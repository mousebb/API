package categoryCtlr

import (
	"github.com/curt-labs/API/middleware"

	"github.com/curt-labs/API/models/category"

	"net/http"
	"strconv"
)

// GetCategory ...
func GetCategory(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	var c category.Category
	var err error

	c.CategoryID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil || c.CategoryID == 0 {
		return nil, err
	}

	qs := r.URL.Query()
	page := 1
	count := 50
	if pg := qs.Get("page"); pg != "" {
		page, _ = strconv.Atoi(pg)
	}
	if ct := qs.Get("count"); ct != "" {
		count, _ = strconv.Atoi(ct)
	}

	err = c.Get(page, count)
	if err != nil || c.CategoryID == 0 {
		return nil, err
	}

	return c, nil
}

func GetCategoryTree(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) ([]category.Category, error) {
	cats, err := category.GetCategoryTree(ctx.DataContext)
	if err != nil {
		return cats, err
	}

	return cats, nil
}

func GetCategoryParts(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (category.PartResponse, error) {
	catId, err := strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return category.PartResponse{}, err
	}

	qs := r.URL.Query()
	page := 1
	count := 50
	if pg := qs.Get("page"); pg != "" {
		page, _ = strconv.Atoi(pg)
	}
	if ct := qs.Get("count"); ct != "" {
		count, _ = strconv.Atoi(ct)
	}

	return category.GetCategoryParts(catId, page, count)
}
