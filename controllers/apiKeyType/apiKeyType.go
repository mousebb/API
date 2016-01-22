package apiKeyType

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/apiKeyType"

	"net/http"
)

// GetAPIKeyTypes Returns a list of available API key types.
func GetAPIKeyTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return nil, err
	}

	return apiKeyType.GetAllKeyTypes(tx)
}
