package categoryCtlr

import (
	"github.com/curt-labs/API/middleware"
	"log"

	"github.com/curt-labs/API/models/category"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"strconv"
)

// GetCategory ...
func GetCategory(ctx *middleware.ApiContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	ps := context.Get(r, "params").(httprouter.Params)

	var c category.Category
	var err error

	log.Println(ps.ByName("id"))
	c.CategoryID, err = strconv.Atoi(ps.ByName("id"))
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

//
// func GetCategoryTree(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) string {
// 	cats, err := category.GetCategoryTree(dtx)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting categories", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(cats))
// }
//
// func GetCategoryParts(rw http.ResponseWriter, r *http.Request, params martini.Params, enc encoding.Encoder, dtx *apicontext.DataContext) string {
// 	catIdStr := params["id"]
// 	catId, err := strconv.Atoi(catIdStr)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting category Id", err, rw, r)
// 		return ""
// 	}
//
// 	qs := r.URL.Query()
// 	page := 1
// 	count := 50
// 	if pg := qs.Get("page"); pg != "" {
// 		page, _ = strconv.Atoi(pg)
// 	}
// 	if ct := qs.Get("count"); ct != "" {
// 		count, _ = strconv.Atoi(ct)
// 	}
//
// 	parts, err := category.GetCategoryParts(catId, page, count)
// 	if err != nil {
// 		apierror.GenerateError("Trouble getting parts", err, rw, r)
// 		return ""
// 	}
//
// 	return encoding.Must(enc.Encode(parts))
// }
