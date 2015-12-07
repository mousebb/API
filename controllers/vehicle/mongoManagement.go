package vehicle

import (
	"fmt"

	"github.com/curt-labs/API/helpers/encoding"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetAllCollectionApplications(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}
	apps, err := products.GetAllCollectionApplications(collection)
	if err != nil {
		apierror.GenerateError(err.Error(), err, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(apps))
}

func UpdateApplication(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var app products.NoSqlVehicle
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Error reading request body", nil, w, r)
		return ""
	}

	if err = json.Unmarshal(body, &app); err != nil {
		apierror.GenerateError("Error decoding vehicle application", nil, w, r)
		return ""
	}

	if err = app.Update(collection); err != nil {
		apierror.GenerateError("Error updating vehicle", nil, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(app))
}

func DeleteApplication(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var app products.NoSqlVehicle
	collection := params["collection"]
	if collection == "" {
		apierror.GenerateError("No Collection in URL", nil, w, r)
		return ""
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		apierror.GenerateError("Error reading request body", nil, w, r)
		return ""
	}

	if err = json.Unmarshal(body, &app); err != nil {
		apierror.GenerateError("Error decoding vehicle application", nil, w, r)
		return ""
	}

	if err = app.Delete(collection); err != nil {
		apierror.GenerateError("Error updating vehicle", nil, w, r)
		return ""
	}
	return encoding.Must(enc.Encode(app))
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

	err = app.Create(collection)

	return app, err
}
