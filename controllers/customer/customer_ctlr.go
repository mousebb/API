package customerCtlr

import (
	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/products"

	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetCustomer Retrieves a customer.Customer based of the given API key.
func GetCustomer(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	var c customer.Customer

	if err = c.GetCustomerIdFromKey(ctx); err != nil {
		return nil, err
	}

	if err = c.GetCustomer(ctx); err != nil {
		return nil, err
	}

	lowerKey := strings.ToLower(ctx.DataContext.APIKey)
	for i, u := range c.Users {
		for _, k := range u.Keys {
			if strings.ToLower(k.Key) == lowerKey {
				c.Users[i].Current = true
			}
		}
	}

	return c, nil
}

// GetLocations Returns the []customer.Location for a given customer.Customer.
func GetLocations(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {

	c, err := ctx.Cu.GetCustomer(ctx.DataContext.APIKey)
	if err != nil {
		return nil, err
	}

	return c.Locations, nil
}

// GetUsers Returns the []customer.User for a given customer.Customer.
func GetUsers(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	user, err := customer.GetCustomerUserFromKey(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer user", err, rw, r)
		return ""
	}

	if !user.Sudo {
		err = errors.New("Unauthorized!")
		apierror.GenerateError("Unauthorized!", err, rw, r, http.StatusUnauthorized)
		return ""
	}

	cust, err := user.GetCustomer(dtx.APIKey)
	if err != nil {
		apierror.GenerateError("Trouble getting customer", err, rw, r)
		return ""
	}

	if err = cust.GetUsers(dtx.APIKey); err != nil {
		apierror.GenerateError("Trouble getting users", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(cust.Users))
}

// GetUser Returns the customer.User for the supplied API key.
func GetUser(rw http.ResponseWriter, r *http.Request, enc encoding.Encoder) string {
	qs, err := url.Parse(r.URL.String())
	if err != nil {
		apierror.GenerateError("err parsing url", err, rw, r)
		return ""
	}

	key := qs.Query().Get("key")
	if key == "" {
		key = r.FormValue("key")
	}

	user, err := customer.GetCustomerUserFromKey(key)
	if err != nil {
		apierror.GenerateError("Trouble getting customer user", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(user))
}

// GetCustomerPrice Returns the selling price of an item for a given customer.Customer.
func GetCustomerPrice(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	var p products.Part

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if p.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	var price float64
	if price, err = customer.GetCustomerPrice(dtx, p.ID); err != nil {
		apierror.GenerateError("Trouble getting price", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(price))
}

// GetCustomerCartReference Returns the cooresponding identifier to the given product
// identifier for the customer.
func GetCustomerCartReference(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var err error
	var p products.Part

	id := r.FormValue("id")
	if id == "" {
		id = params["id"]
	}

	if p.ID, err = strconv.Atoi(id); err != nil {
		apierror.GenerateError("Trouble getting part ID", err, rw, r)
		return ""
	}

	var ref int
	if ref, err = customer.GetCustomerCartReference(dtx.APIKey, p.ID); err != nil {
		apierror.GenerateError("Trouble getting customer cart reference", err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(ref))
}

// SaveCustomer Updates the supplied customer.Customer.
func SaveCustomer(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	var c customer.Customer
	var err error

	if r.FormValue("id") != "" || params["id"] != "" {
		id := r.FormValue("id")
		if id == "" {
			id = params["id"]
		}

		if c.Id, err = strconv.Atoi(id); err != nil {
			apierror.GenerateError("Trouble getting customer ID", err, rw, r)
			return ""
		}

		if err = c.Basics(dtx.APIKey); err != nil {
			apierror.GenerateError("Trouble getting customer", err, rw, r)
			return ""
		}
	}

	//json
	var requestBody []byte
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		apierror.GenerateError("Trouble reading request body while saving customer", err, rw, r)
		return ""
	}

	if err = json.Unmarshal(requestBody, &c); err != nil {
		apierror.GenerateError("Trouble unmarshalling json request body while saving customer", err, rw, r)
		return ""
	}

	//create or update
	if c.Id > 0 {
		err = c.Update()
	} else {
		err = c.Create()
	}

	if err != nil {
		msg := "Trouble creating customer"
		if c.Id > 0 {
			msg = "Trouble updating customer"
		}
		apierror.GenerateError(msg, err, rw, r)
		return ""
	}

	return encoding.Must(enc.Encode(c))
}
