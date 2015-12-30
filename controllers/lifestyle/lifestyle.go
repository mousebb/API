package lifestyle

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/lifestyle"

	"net/http"
	"strconv"
)

func GetAll(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	return lifestyle.GetAll(ctx)
}

func Get(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var l lifestyle.Lifestyle
	var err error
	l.ID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return nil, err
	}

	err = l.Get(ctx)

	return l, err
}
