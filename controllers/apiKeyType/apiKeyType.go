package apiKeyType

import (
	"github.com/curt-labs/API/models/apiKeyType"
	"github.com/curt-labs/API/middleware"

	"net/http"
)

func GetApiKeyTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return apiKeyType.GetAllApiKeyTypes()
}
