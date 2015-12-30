package vinLookup

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/vinLookup"
)

func GetParts(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	vin := ctx.Params.ByName("vin")

	return vinLookup.VinPartLookup(ctx, vin)
}

func GetConfigs(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	vin := ctx.Params.ByName("vin")

	return vinLookup.GetVehicleConfigs(ctx, vin)
}

func GetPartsFromVehicleID(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	vehicleID := ctx.Params.ByName("vehicleID")
	id, err := strconv.Atoi(vehicleID)
	if err != nil {
		return nil, err
	}

	v := vinLookup.CurtVehicle{
		ID: id,
	}

	return v.GetPartsFromVehicleConfig(ctx)
}
