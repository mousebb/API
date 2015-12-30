package showcase

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/showcase"
)

func GetAllShowcases(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
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

	return showcase.GetAllShowcases(ctx, page, count, randomize)
}

func GetShowcase(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var show showcase.Showcase

	if show.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = show.Get(ctx)

	return show, err
}

func Save(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var show showcase.Showcase
	var err error
	idStr := ctx.Params.ByName("id")
	if idStr != "" {
		show.ID, err = strconv.Atoi(idStr)
		err = show.Get(ctx)
		if err != nil {
			return nil, err
		}
	}
	//json
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(requestBody, &show)
	if err != nil {
		return nil, err
	}
	//create or update
	if show.ID > 0 {
		err = show.Update(ctx)
	} else {
		err = show.Create(ctx)
	}

	return show, err
}

func Delete(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var a showcase.Showcase

	idStr := ctx.Params.ByName("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	a.ID = id
	err = a.Delete(ctx)

	return a, err
}
