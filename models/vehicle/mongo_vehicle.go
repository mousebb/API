package vehicle

import (
	"log"

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

func GetYears(ctx *middleware.APIContext) ([]string, error) {
	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$unwind": "$vehicle_applications",
		},
		bson.M{
			"$group": bson.M{"_id": "$vehicle_applications.year"},
		},
		bson.M{
			"$sort": bson.M{"_id": -1},
		},
	}

	type YearResult struct {
		Year string `bson:"_id"`
	}

	var res []YearResult
	err := c.Pipe(aggr).All(&res)
	if err != nil {
		return nil, err
	}

	var yrs []string
	for _, y := range res {
		yrs = append(yrs, y.Year)
	}

	return yrs, err
}

func GetMakes(ctx *middleware.APIContext, year string) ([]string, error) {
	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$match": bson.M{"vehicle_applications.year": year},
		},
		bson.M{
			"$unwind": "$vehicle_applications",
		},
		bson.M{
			"$group": bson.M{"_id": "$vehicle_applications.make"},
		},
		bson.M{
			"$sort": bson.M{"_id": 1},
		},
	}

	type MakeResult struct {
		Make string `bson:"_id"`
	}

	var res []MakeResult
	err := c.Pipe(aggr).All(&res)
	if err != nil {
		return nil, err
	}

	var yrs []string
	for _, y := range res {
		yrs = append(yrs, y.Make)
	}

	return yrs, err
}

func GetModels(ctx *middleware.APIContext, year, make string) ([]string, error) {
	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year": year,
				"vehicle_applications.make": make,
			},
		},
		bson.M{
			"$unwind": "$vehicle_applications",
		},
		bson.M{
			"$group": bson.M{"_id": "$vehicle_applications.model"},
		},
		bson.M{
			"$sort": bson.M{"_id": 1},
		},
	}
	log.Println(aggr)

	type ModelResult struct {
		Model string `bson:"_id"`
	}

	var res []ModelResult
	err := c.Pipe(aggr).All(&res)
	if err != nil {
		return nil, err
	}

	var yrs []string
	for _, y := range res {
		yrs = append(yrs, y.Model)
	}

	return yrs, err
}
