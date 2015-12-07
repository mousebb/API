package search_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/search"
)

func Search(ctx *middleware.APIConext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	terms := ctx.Params.ByName("term")
	qs := r.URL.Query()
	page, _ := strconv.Atoi(qs.Get("page"))
	count, _ := strconv.Atoi(qs.Get("count"))
	brand, _ := strconv.Atoi(qs.Get("brand"))

	return search.Dsl(ctx, terms, page, count, brand)
}
