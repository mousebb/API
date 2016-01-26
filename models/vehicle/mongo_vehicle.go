package vehicle

import (
	"fmt"
	"sort"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"gopkg.in/mgo.v2/bson"
)

type VehicleApplication struct {
	Year        string `bson:"year" json:"year" xml:"year"`
	Make        string `bson:"make" json:"make" xml:"make"`
	Model       string `bson:"model" json:"model" xml:"model"`
	Style       string `bson:"style" json:"style" xml:"style"`
	Exposed     string `bson:"exposed" json:"exposed" xml:"exposed"`
	Drilling    string `bson:"drilling" json:"drilling" xml:"drilling"`
	InstallTime string `bson:"install_time" json:"install_time" xml:"install_time"`
}

func ReverseMongoLookup(ctx *middleware.APIContext, part string) ([]VehicleApplication, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"part_number": part,
	}

	var apps []VehicleApplication
	err := c.Find(qry).Select(bson.M{"vehicle_applications": 1, "_id": 0}).All(&apps)

	return apps, err
}

func GetYears(ctx *middleware.APIContext) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": []int{800, 900},
		},
	}

	var res []string
	err := c.Find(qry).Select(bson.M{
		"vehicle_applications.year": 1,
		"_id": -1,
	}).Distinct("vehicle_applications.year", &res)

	if err != nil {
		return nil, err
	}

	sort.Sort(sort.Reverse(sort.StringSlice(res)))

	return res, err
}

func GetMakes(ctx *middleware.APIContext, year string) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year": year,
			},
		},
		bson.M{
			"$unwind": "$vehicle_applications",
		},
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year": year,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$vehicle_applications.make",
			},
		},
	}

	type Result struct {
		Make string `bson:"_id"`
	}

	var res []Result
	err := c.Pipe(aggr).All(&res)

	if err != nil {
		return nil, err
	}

	var makes []string
	existing := make(map[string]string, 0)
	for _, r := range res {
		if _, ok := existing[r.Make]; !ok {
			makes = append(makes, r.Make)
			existing[r.Make] = r.Make
		}
	}

	sort.Strings(makes)

	return makes, err
}

func GetModels(ctx *middleware.APIContext, year, vehicleMake string) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year": year,
				"vehicle_applications.make": vehicleMake,
			},
		},
		bson.M{
			"$unwind": "$vehicle_applications",
		},
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year": year,
				"vehicle_applications.make": vehicleMake,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$vehicle_applications.model",
			},
		},
	}

	type Result struct {
		Model string `bson:"_id"`
	}

	var res []Result
	err := c.Pipe(aggr).All(&res)

	if err != nil {
		return nil, err
	}

	var models []string
	existing := make(map[string]string, 0)
	for _, r := range res {
		if _, ok := existing[r.Model]; !ok {
			models = append(models, r.Model)
			existing[r.Model] = r.Model
		}
	}

	sort.Strings(models)

	return models, err
}

func GetStyles(ctx *middleware.APIContext, year, vehicleMake, model string) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	aggr := []bson.M{
		bson.M{
			"$match": bson.M{
				"vehicle_applications.year":  year,
				"vehicle_applications.make":  vehicleMake,
				"vehicle_applications.model": model,
			},
		},
		bson.M{
			"$project": bson.M{
				"styles": bson.M{
					"$filter": bson.M{
						"input": "$vehicle_applications",
						"as":    "apps",
						"cond": bson.M{
							"$and": []bson.M{
								bson.M{"$eq": []string{"$$apps.year", year}},
								bson.M{"$eq": []string{"$$apps.make", vehicleMake}},
								bson.M{"$eq": []string{"$$apps.model", model}},
							},
						},
					},
				},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": "$styles.style",
			},
		},
	}

	type Result struct {
		Styles []string `bson:"_id"`
	}

	var res []Result
	err := c.Pipe(aggr).All(&res)

	if err != nil {
		return nil, err
	}

	var styles []string
	existing := make(map[string]string, 0)
	for _, r := range res {
		for _, s := range r.Styles {
			if _, ok := existing[s]; !ok {
				styles = append(styles, s)
				existing[s] = s
			}
		}
	}

	sort.Strings(styles)

	return styles, err
}
