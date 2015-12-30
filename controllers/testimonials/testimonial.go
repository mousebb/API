package testimonials

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/testimonials"
)

func GetAllTestimonials(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var page int
	var count int
	var randomize bool

	qs := req.URL.Query()

	if qs.Get("page") != "" {
		if pg, err := strconv.Atoi(qs.Get("page")); err == nil {
			page = pg
		}
	}
	if qs.Get("count") != "" {
		if c, err := strconv.Atoi(qs.Get("count")); err == nil {
			count = c
		}
	}

	if qs.Get("randomize") != "" {
		randomize, _ = strconv.ParseBool(qs.Get("randomize"))
	}

	return testimonials.GetAllTestimonials(ctx, page, count, randomize)
}

func GetTestimonial(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var test testimonials.Testimonial

	if test.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}
	err = test.Get(ctx)

	return test, err
}

func Save(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var a testimonials.Testimonial
	var err error
	idStr := ctx.Params.ByName("id")
	if idStr != "" {
		a.ID, err = strconv.Atoi(idStr)
		err = a.Get(ctx)
		if err != nil {
			return nil, err
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(requestBody, &a)
	if err != nil {
		return nil, err
	}
	//create or update
	if a.ID > 0 {
		err = a.Update(ctx)
	} else {
		err = a.Create(ctx)
	}

	return a, err
}

func Delete(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var a testimonials.Testimonial

	idStr := ctx.Params.ByName("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	a.ID = id
	err = a.Delete(ctx)

	return a, err
}
