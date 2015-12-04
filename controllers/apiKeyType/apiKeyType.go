package apiKeyType

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/apiKeyType"

	"net/http"
)

// GetApiKeyTypes Returns a list of available API key types.
func GetApiKeyTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return apiKeyType.GetAllApiKeyTypes(ctx)
}
