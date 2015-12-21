package middleware

import (
	"fmt"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/iapi/brand"
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
	var dtx DataContext
	var u customer.User

	c := ctx.Session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	err := c.Find(bson.M{"keys.$.key": k, "keys.$.key.type.type": t, "active": 1}).One(&u)
	if err != nil {
		return nil, err
	}
	if u.ID == "" {
		return nil, fmt.Errorf("failed to find account for the provided %s key: %s", t, k)
	}

	dtx.brandArray()
	dtx.brandString()

	dtx.User = u
	dtx.APIKey = k

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
