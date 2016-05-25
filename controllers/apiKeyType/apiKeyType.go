// swagger generator: swagger -apiPackage="github.com/curt-labs/API/controllers/apiKeyType" -mainApiFile="github.com/curt-labs/API/index.go"

// Package apiKeyType Allows request operations for API Key Types.
// @SubApi API Key Management [/api]
package apiKeyType

import (
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/apiKeyType"

	"net/http"
)

// GetAPIKeyTypes Returns all availabe APIKeyTypes
// @Title GetAPIKeyTypes
// @Description Returns all the available types of API keys
// @Accept  json
// @Param   key     query    string     true        "Public API Key"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Resource /api
// @Router /api/keys/types [get]
func GetAPIKeyTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	tx, err := ctx.DB.Begin()
	if err != nil {
		return nil, err
	}

	return apiKeyType.GetAllKeyTypes(tx)
}
