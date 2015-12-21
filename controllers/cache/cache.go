package cache

import (
	"fmt"
	"net/http"

	"github.com/curt-labs/API/helpers/redis"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
)

//TODO check for super user

// GetKeys ...
func GetKeys(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	if !approveuser(ctx) {
		return nil, fmt.Errorf("%s", "unauthorized request")
	}

	return redis.GetNamespaces()
}

// GetByKey ...
func GetByKey(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	if !approveuser(ctx) {
		return nil, fmt.Errorf("%s", "unauthorized request")
	}

	key := r.URL.Query().Get("redis_key")
	namespace := r.URL.Query().Get("redis_namespace")

	return redis.GetFullPath(fmt.Sprintf("%s:%s", namespace, key))
}

// DeleteKey ...
func DeleteKey(ctx *middleware.APIContext, rw http.ResponseWriter, r *http.Request) (interface{}, error) {
	if !approveuser(ctx) {
		return nil, fmt.Errorf("%s", "unauthorized request")
	}

	key := r.URL.Query().Get("redis")

	return nil, redis.DeleteFullPath(key)
}

func approveuser(ctx *middleware.APIContext) bool {

	c := customer.Customer{}
	var err error
	// TODO - re-enable this
	// err = c.GetCustomerIdFromKey(ctx)
	if err != nil || c.ID == 0 {
		return false
	}
	return true
}
