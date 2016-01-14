package customerCtlr

import (
	"net/http"
	"strings"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
)

// GetUser This can only be called by supplying the `Private` APIKey
// for the User that is being requested. Also, this endpoint requires
// `sudo` privileges on the User related to the requestor's APIKey (?key=).
func GetUser(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return customer.GetUserByKey(ctx.Session, ctx.Params.ByName("key"), "Private")
}

// GetUserByKey This can only be called by supplying the `Private` APIKey
// in the request query string (?key=).
func GetUserByKey(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return customer.GetUserByKey(ctx.Session, ctx.DataContext.APIKey, "Private")
}

// Authenticate Will take in a username/password combination and
// authenticate a User, returning the information associated to that User.
func Authenticate(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var email string
	var pass string

	email = strings.TrimSpace(r.FormValue("email"))
	pass = strings.TrimSpace(r.FormValue("password"))

	return customer.AuthenticateUser(ctx.Session, email, pass)
}
