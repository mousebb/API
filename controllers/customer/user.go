package customerCtlr

import (
	"net/http"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
)

// GetUser This can only be called by supplying the `Private` APIKey
// for the User that is being requested. Also, this endpoint requires
// `sudo` privileges on the User related to the requestor's APIKey (?key=).
func GetUser(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return customer.AuthenticateUserByKey(ctx.Session, ctx.Params.ByName("key"))
}
