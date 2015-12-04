package vehicle

import (
	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"gopkg.in/mgo.v2/bson"
)

type MgoVehicle struct {
	Identifier bson.ObjectId `bson:"_id" json:"-" xml:"-"`
	Year       string        `bson:"year" json:"year,omitempty" xml:"year,omitempty"`
	Make       string        `bson:"make" json:"make,omitempty" xml:"make,omitempty"`
	Model      string        `bson:"model" json:"model,omitempty" xml:"model,omitempty"`
	Style      string        `bson:"style" json:"style,omitempty" xml:"style,omitempty"`
}

const (
	AriesDb = "aries"
)

func ReverseMongoLookup(partId int, ctx *middleware.APIContext) (vehicles []MgoVehicle, err error) {

	collections, err := ctx.AriesSession.DB(database.AriesMongoDatabase).CollectionNames()
	if err != nil {
		return
	}
	for _, collection := range collections {
		var temps []MgoVehicle
		query := bson.M{
			"parts": partId,
		}
		err = ctx.AriesSession.DB(database.AriesMongoDatabase).C(collection).Find(query).All(&temps)
		if err != nil {
			return
		}
		vehicles = append(vehicles, temps...)
	}
	return
}
