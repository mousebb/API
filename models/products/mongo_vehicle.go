package products

import (
	"fmt"
	"sort"
	"strings"

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
	InstallTime string `bson:"install_time" json:"installTime" xml:"installTime"`

	Years  []string `json:"availableYears,omitempty" xml:"availableYears,omitempty"`
	Makes  []string `json:"availableMakes,omitempty" xml:"availableMakes,omitempty"`
	Models []string `json:"availableModels,omitempty" xml:"availableModels,omitempty"`

	CategoryStyles []CategoryStyle `json:"categoryStyles" xml:"categoryStyles"`
}

type CategoryStyle struct {
	Category   Category     `json:"category" xml:"category"`
	StyleParts []StyleParts `json:"styleParts" xml:"styleParts"`
}

type StyleParts struct {
	Style string `json:"style" xml:"style"`
	Parts []Part `json:"parts" xml:"parts"`
}

func (va *VehicleApplication) Query(ctx *middleware.APIContext) error {
	var err error

	if va.Year == "" {
		va.Years, err = getYears(ctx)
	} else if va.Year != "" && va.Make == "" {
		va.Makes, err = getMakes(ctx, va.Year)
	} else if va.Year != "" && va.Make != "" && va.Model == "" {
		va.Models, err = getModels(ctx, va.Year, va.Make)
	} else if va.Year != "" && va.Make != "" && va.Model != "" {
		va.CategoryStyles, err = getStyles(ctx, va.Year, va.Make, va.Model)
	}

	return err
}

func ReverseMongoLookup(ctx *middleware.APIContext, part string) ([]VehicleApplication, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	} else if ctx.DataContext == nil {
		return nil, fmt.Errorf("invalid data context")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"part_number": part,
	}

	var apps []VehicleApplication
	err := c.Find(qry).Select(bson.M{"vehicle_applications": 1, "_id": 0}).All(&apps)

	return apps, err
}

func getYears(ctx *middleware.APIContext) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	} else if ctx.DataContext == nil {
		return nil, fmt.Errorf("invalid data context")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	qry := bson.M{
		"status": bson.M{
			"$in": []int{800, 900},
		},
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
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

func getMakes(ctx *middleware.APIContext, year string) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	} else if ctx.DataContext == nil {
		return nil, fmt.Errorf("invalid data context")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
			},
		},
		"status": bson.M{
			"$in": []int{800, 900},
		},
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}).Select(bson.M{"vehicle_applications.$.make": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var makes []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Make = strings.Title(a.Make)
			if _, ok := existing[a.Make]; !ok {
				makes = append(makes, a.Make)
				existing[a.Make] = a.Make
			}
		}
	}

	sort.Strings(makes)

	return makes, err
}

func getModels(ctx *middleware.APIContext, year, vehicleMake string) ([]string, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	} else if ctx.DataContext == nil {
		return nil, fmt.Errorf("invalid data context")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps []VehicleApplication `bson:"vehicle_applications"`
	}

	var apps []Apps
	err := c.Find(bson.M{
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": []int{800, 900},
		},
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}).Select(bson.M{"vehicle_applications.$.model": 1, "_id": 0}).All(&apps)
	if err != nil {
		return nil, err
	}

	var models []string

	existing := make(map[string]string, 0)
	for _, app := range apps {
		for _, a := range app.Apps {
			a.Model = strings.Title(a.Model)
			if _, ok := existing[a.Model]; !ok {
				models = append(models, a.Model)
				existing[a.Model] = a.Model
			}
		}
	}

	sort.Strings(models)

	return models, err
}

func getStyles(ctx *middleware.APIContext, year, vehicleMake, model string) ([]CategoryStyle, error) {
	if ctx.Session == nil {
		return nil, fmt.Errorf("invalid mongodb connection")
	} else if ctx.DataContext == nil {
		return nil, fmt.Errorf("invalid data context")
	}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName)

	type Apps struct {
		Apps    []VehicleApplication `bson:"vehicle_applications"`
		PartNum string               `bson:"part_number"`
	}

	var parts []Part
	err := c.Find(bson.M{
		"vehicle_applications": bson.M{
			"$elemMatch": bson.M{
				"year": year,
				"make": bson.RegEx{
					Pattern: "^" + vehicleMake + "$",
					Options: "i",
				},
				"model": bson.RegEx{
					Pattern: "^" + model + "$",
					Options: "i",
				},
			},
		},
		"status": bson.M{
			"$in": []int{800, 900},
		},
		"brand.id": bson.M{
			"$in": ctx.DataContext.BrandArray,
		},
	}).All(&parts)

	if err != nil {
		return nil, err
	}

	lowerMake := strings.ToLower(vehicleMake)
	lowerModel := strings.ToLower(model)

	css := make(map[string]CategoryStyle, 0)
	for _, p := range parts {
		if len(p.Categories) == 0 {
			continue
		}

		for _, va := range p.Vehicles {
			if va.Year != year || strings.ToLower(va.Make) != lowerMake || strings.ToLower(va.Model) != lowerModel {
				continue
			}

			css = mapPartToCategoryStyles(p, css, va.Style)
		}

	}

	var catStyles []CategoryStyle
	for _, cs := range css {
		catStyles = append(catStyles, cs)
	}

	return catStyles, nil
}

func getChildCategory(cats []Category) (Category, error) {
	for _, cat := range cats {
		if len(cat.Children) == 0 {
			return cat, nil
		}
	}

	return Category{}, fmt.Errorf("failed to locate child")
}

func mapPartToCategoryStyles(p Part, css map[string]CategoryStyle, style string) map[string]CategoryStyle {
	childCat, err := getChildCategory(p.Categories)
	if err != nil || childCat.Identifier.String() == "" {
		return css
	}

	cs, ok := css[childCat.Identifier.String()]
	if !ok {
		cs = CategoryStyle{
			Category: childCat,
		}
	}

	for i, sp := range cs.StyleParts {
		if strings.ToLower(sp.Style) == strings.ToLower(style) {
			cs.StyleParts[i].Parts = append(cs.StyleParts[i].Parts, p)
			css[childCat.Identifier.String()] = cs
			return css
		}
	}

	currentStyle := StyleParts{
		Style: strings.Title(style),
		Parts: []Part{p},
	}

	cs.StyleParts = append(cs.StyleParts, currentStyle)
	css[childCat.Identifier.String()] = cs

	return css
}
