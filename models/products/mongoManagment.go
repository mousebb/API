package products

import (
	"github.com/curt-labs/API/middleware"
	"gopkg.in/mgo.v2/bson"
)

func GetAllCollectionApplications(ctx *middleware.APIContext, collection string) ([]NoSqlVehicle, error) {
	var apps []NoSqlVehicle

	err := ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(collection).Find(bson.M{}).Sort("-year", "make", "model", "style").All(&apps)

	return apps, err
}

func (n *NoSqlVehicle) Update(ctx *middleware.APIContext, collection string) error {

	update := make(map[string]interface{})
	if n.Year != "" {
		update["year"] = n.Year
	}
	if n.Make != "" {
		update["make"] = n.Make
	}
	if n.Model != "" {
		update["model"] = n.Model
	}
	if n.Style != "" {
		update["style"] = n.Style
	}
	if n.Make != "" {
		update["make"] = n.Make
	}
	if len(n.PartIdentifiers) > 0 {
		update["parts"] = n.PartIdentifiers
	}
	return ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(collection).UpdateId(n.ID, update)
}

func (n *NoSqlVehicle) Delete(ctx *middleware.APIContext, collection string) error {
	return ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(collection).RemoveId(n.ID)
}

func (n *NoSqlVehicle) Create(ctx *middleware.APIContext, collection string) error {
	return ctx.AriesSession.DB(ctx.AriesMongoDatabase).C(collection).Insert(n)
}
