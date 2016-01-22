package customerCtlr

import (
	"encoding/json"
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

// GetUserByIdentifier Returns a User by retrieving based off the provided ``:id`
// parameter. This can only be called by a User with super privileges.
func GetUserByIdentifier(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	return customer.GetUser(ctx.Session, ctx.Params.ByName("id"), ctx.DataContext.APIKey)
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

// AddUser Will commit a new user to the same Customer object as
// the requestor's Customer reference. It will not update the following
// fields from the submitted User object: `ID`, `CustomerNumber`, `DateAdded`, `Keys`, or `ComnetAccounts`.
func AddUser(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user *customer.User

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	err = customer.AddUser(ctx.Session, ctx.DB, user, ctx.DataContext.APIKey)

	return user, err
}

// UpdateUser Can update the Name, Email, SuperUser (if updated by a super user).
// If the update is called by a different requestor than the updating User, the
// requestor is required to be have super powers.
func UpdateUser(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var user *customer.User
	var err error
	superable := false

	user.ID = ctx.Params.ByName("id")

	if user.ID != "" {
		superable = true
		user, err = customer.GetUser(ctx.Session, user.ID, ctx.DataContext.APIKey)
	} else {
		user, err = customer.GetUserByKey(ctx.Session, ctx.DataContext.APIKey, customer.PrivateKeyType)
	}

	if err != nil {
		return nil, err
	}

	var changeUser customer.User
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&changeUser)
	if err != nil {
		return nil, err
	}

	if changeUser.Location != nil {
		user.Location = changeUser.Location
	}

	user.Name = changeUser.Name
	user.Email = changeUser.Email
	if superable {
		user.SuperUser = changeUser.SuperUser
	}

	err = customer.UpdateUser(ctx.DB, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
