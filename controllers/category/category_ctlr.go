package category_ctlr

import (
	"github.com/curt-labs/API/helpers/apicontext"
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/models/category"
	"github.com/go-martini/martini"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"

	"net/http"
	"strconv"
)

// GetCategory
func GetCategory(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var c category.Category
	var err error
	c.CategoryID, err = strconv.Atoi(p.ByName("id"))
	if err != nil || c.CategoryID == 0 {
		apierror.GenerateError("Trouble getting category identifier", err, rw, r)
		return
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
		apierror.GenerateError("Trouble getting category", err, rw, r)
		return
	}

	rndr, ok := context.GetOk(r, "renderer")
	if !ok {
		apierror.GenerateError("Trouble getting category", err, rw, r)
		return
	}

	rndr.(render.Render).JSON(rw, 200, c)

	return
}

func GetCategoryTree(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	ctx := context.GetOk(r, key)
	cats, err := category.GetCategoryTree(dtx)
	if err != nil {
		apierror.GenerateError("Trouble getting categories", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cats))
}

func GetCategoryParts(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
	catIdStr := params["id"]
	catId, err := strconv.Atoi(catIdStr)
	if err != nil {
		apierror.GenerateError("Trouble getting category Id", err, rw, r)
		return ""
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

	parts, err := category.GetCategoryParts(catId, page, count)
	if err != nil {
		apierror.GenerateError("Trouble getting parts", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(parts))
}
