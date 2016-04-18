package customer

import (
	"fmt"
	"strings"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/brand"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DataContext Holds the relevant data
// that is generated in the middleware
// associated with the current API request.
type DataContext struct {
	BrandID     int
	WebsiteID   int
	APIKey      string
	CustomerID  int
	User        User
	Brands      []brand.Brand
	BrandString string
	BrandArray  []int
}

// NewContext Retrieves the required information from MongoDB associated
// with the API key and type provided.
func NewContext(sess *mgo.Session, k string, t string, requireSudo bool) (*DataContext, error) {
	resp := struct {
		Users []User
	}{}

	c := sess.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName)

	// TODO: This needs to be change to allow non-active users to
	// be aggregated out. However, they're currently not being
	// fanned into MongoDB.
	//
	// we don't need the `active` operator, since only active users
	// are being put into MongoDB.
	qry := bson.M{"users.keys.key": k}
	if t != "" {
		pattern := bson.RegEx{
			Pattern: "^" + t + "$",
			Options: "i",
		}
		qry["users.keys.type.type"] = pattern
	}
	if requireSudo {
		qry["users.superUser"] = true
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
