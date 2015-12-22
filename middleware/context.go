package middleware

import (
	"fmt"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DataContext struct {
	BrandID     int
	WebsiteID   int
	APIKey      string
	CustomerID  int
	User        customer.User
	Brands      []brand.Brand
	BrandString string
	BrandArray  []int
}

func (ctx *APIContext) BuildDataContext(k, t string) error {
	dtx, err := new(ctx, k, t)
	if err != nil {
		return err
	}

	ctx.DataContext = dtx

	return nil
}

func new(ctx *APIContext, k, t string) (*DataContext, error) {
	resp := struct {
		Users []customer.User
	}{}

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	// we don't need the `active` operator, since only active users
	// are being put into MongoDB.
	qry := bson.M{"users.keys.key": k}
	if t != "" {
		qry["users.keys.type.type"] = t
	}

	err := c.Find(qry).Select(bson.M{"users.$.user": 1, "_id": 0}).One(&resp)
	if err != nil && err != mgo.ErrNotFound {
		return nil, err
	}

	if len(resp.Users) == 0 {
		return nil, fmt.Errorf("failed to find account for the provided %s key: %s", t, k)
	}

	dtx := DataContext{
		User:   resp.Users[0],
		APIKey: k,
	}

	for _, apiKey := range dtx.User.Keys {
		if apiKey.Key == k {
			dtx.Brands = apiKey.Brands
		}
	}

	dtx.brandArray()
	dtx.brandString()

	return &dtx, nil
}

func (dtx *DataContext) brandString() {
	var ids []string
	for _, b := range dtx.Brands {
		ids = append(ids, string(b.ID))
	}
	dtx.BrandString = strings.Join(ids, ",")
}

func (dtx *DataContext) brandArray() {
	dtx.BrandArray = []int{}
	for _, b := range dtx.Brands {
		dtx.BrandArray = append(dtx.BrandArray, b.ID)
	}
}
