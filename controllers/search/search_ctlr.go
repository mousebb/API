package searchCtlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/search"
)

// Search Uses URL paramter :term to query elastic search.
// Allows paging via query string parameters `page` and `count`.
// Brand designation is also available via query string `brand`.
func Search(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	terms := ctx.Params.ByName("term")
	qs := r.URL.Query()
	page, _ := strconv.Atoi(qs.Get("page"))
	count, _ := strconv.Atoi(qs.Get("count"))
	brand, _ := strconv.Atoi(qs.Get("brand"))

	return search.Dsl(ctx, terms, page, count, brand)
}
