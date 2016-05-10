package customerCtlr

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/models/apiKeyType"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"github.com/ory-am/dockertest"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	db             *sql.DB
	dbName         = "product_data"
	collectionName = "customer"
	session        *mgo.Session

	testCustUser = customer.Customer{
		Identifier: bson.NewObjectId(),
		Users: []customer.User{{
			ID:             "1",
			CustomerNumber: 0,
			Name:           "Test Customer Name",
			Email:          "test@curtmfg.com",
			Keys: []customer.APIKey{{
				Key: "123",
				Brands: []brand.Brand{
					brand.Brand{
						ID: 3,
					},
				},
				Type: apiKeyType.KeyType{
					Type: "Private",
				},
			}},
		}},
	}
	testUsers = []customer.Customer{testCustUser}
)

func TestMain(m *testing.M) {
	var err error

	if os.Getenv("CI") == "" {
		var mongo dockertest.ContainerID
		mongo, err = dockertest.ConnectToMongoDB(3, time.Second*30, func(url string) bool {
			session, err = mgo.Dial(url)
			if err != nil {
				return false
			}
			for _, user := range testUsers {
				err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Insert(user)
				if err != nil {
					log.Fatal(err)
				}
			}
			var u []customer.User
			err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Find(bson.M{}).All(&u)
			return true
		})

		defer func() {
			mongo.KillRemove()
		}()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		info := mgo.DialInfo{
			Addrs:    []string{os.Getenv("MONGO_PORT_27017_TCP_ADDR")},
			Database: "mydb",
			Timeout:  time.Second * 2,
			FailFast: true,
		}

		session, err = mgo.DialWithInfo(&info)
		if err != nil {
			log.Fatal(err)
		}
		for _, user := range testUsers {
			err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Insert(user)
			if err != nil {
				log.Fatal(err)
			}
		}
		var u []customer.User
		err = session.DB(database.ProductMongoDatabase).C(database.CustomerCollectionName).Find(bson.M{}).All(&u)
		if err != nil {
			log.Fatal(err)
		}
	}
	m.Run()
}

// func TestGetAllBrands(t *testing.T) {
// 	Convey("Testing GetUserByKey", t, func() {
// 		ctx := &middleware.APIContext{
// 			DataContext: &customer.DataContext{
// 				BrandID: 3,
// 			},
// 			Params:  httprouter.Params{},
// 			Session: session,
// 		}
//
// 		Convey("with valid db connection", func() {
// 			rec := httptest.NewRecorder()
// 			req, err := http.NewRequest("GET", "http://localhost:8080/customer/user", nil)
// 			So(err, ShouldBeNil)
//
// 			resp, err := GetUserByKey(ctx, rec, req)
// 			So(err, ShouldBeNil)
// 			So(resp, ShouldNotBeNil)
//
// 		})
//
// 	})
// }
