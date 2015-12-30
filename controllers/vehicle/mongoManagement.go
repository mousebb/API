package vehicle

import (
	"fmt"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAllCollectionApplications(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {

	collection := ctx.Params.ByName("collection")
	if collection == "" {
		return nil, fmt.Errorf("collection was empty")
	}

	return products.GetAllCollectionApplications(ctx, collection)
}

func UpdateApplication(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var app products.NoSqlVehicle
	collection := ctx.Params.ByName("collection")
	if collection == "" {
		return nil, fmt.Errorf("collection was empty")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	if err = json.Unmarshal(body, &app); err != nil {
		return nil, fmt.Errorf("error decode vehicle object: %v", err)
	}

	if err = app.Update(ctx, collection); err != nil {
		return nil, fmt.Errorf("error updating vehicle: %v", err)
	}

	return app, nil
}

func DeleteApplication(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var app products.NoSqlVehicle
	collection := ctx.Params.ByName("collection")
	if collection == "" {
		return nil, fmt.Errorf("collection was empty")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	if err = json.Unmarshal(body, &app); err != nil {
		return nil, fmt.Errorf("error decode vehicle object: %v", err)
	}

	if err = app.Delete(ctx, collection); err != nil {
		return nil, fmt.Errorf("error deleting vehicle: %v", err)
	}

	return app, nil
}

func CreateApplication(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var app products.NoSqlVehicle
	collection := ctx.Params.ByName("collection")
	if collection == "" {
		return nil, fmt.Errorf("%s", "a collection must be specified")
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &app); err != nil {
		return nil, err
	}

	err = app.Create(ctx, collection)

	return app, err
}
