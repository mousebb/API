package products

import (
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/brand"
	"github.com/curt-labs/API/models/customer"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	session  *mgo.Session
	emptyCtx *middleware.APIContext
)

func setupMongo() {

	p := Part{
		Identifier: bson.NewObjectId(),
		ID:         12345,
		SKU:        "12345",
		Status:     800,
		Brand: brand.Brand{
			ID:   1,
			Name: "Brand",
			Code: "BA",
		},
		Categories: []Category{
			Category{
				Identifier: bson.NewObjectId(),
				Title:      "Test Category",
			},
		},
		Vehicles: []VehicleApplication{
			VehicleApplication{
				Year:  "2010",
				Make:  "Ford",
				Model: "Fusion",
				Style: "All",
			},
			VehicleApplication{
				Year:  "2009",
				Make:  "Ford",
				Model: "Escape",
				Style: "All, Except Hybrid",
			},
			VehicleApplication{
				Year:  "1997",
				Make:  "Jeep",
				Model: "Grand Cherokee",
				Style: "Laredo",
			},
			VehicleApplication{
				Year:  "1997",
				Make:  "Jeep",
				Model: "Grand Cherokee",
				Style: "Outlander",
			},
			VehicleApplication{
				Year:  "1997",
				Make:  "Ford",
				Model: "F-150",
				Style: "FX4",
			},
		},
	}

	err := session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Insert(&p)
	if err != nil {
		log.Fatal(err)
	}

	p = Part{
		Identifier: bson.NewObjectId(),
		ID:         54321,
		SKU:        "54321",
		Status:     800,
		Brand: brand.Brand{
			ID:   1,
			Name: "Brand",
			Code: "BA",
		},
		Categories: []Category{
			Category{
				Identifier: bson.NewObjectId(),
				Title:      "Test Category 2",
			},
		},
		Vehicles: []VehicleApplication{
			VehicleApplication{
				Year:  "1997",
				Make:  "Jeep",
				Model: "Grand Cherokee",
				Style: "Outlander",
			},
			VehicleApplication{
				Year:  "1997",
				Make:  "Ford",
				Model: "F-150",
				Style: "FX4",
			},
		},
	}

	err = session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Insert(&p)
	if err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	var err error
	if os.Getenv("CI") == "" {
		var mongo dockertest.ContainerID

		mongo, err = dockertest.ConnectToMongoDB(15, time.Second, func(url string) bool {
			session, err = mgo.Dial(url)
			if err != nil {
				log.Fatalf("MongoDB connection failed, with address '%s'.", url)
			}

			session.SetMode(mgo.Monotonic, true)

			setupMongo()

			return session.Ping() == nil
		})

		defer func() {
			session.Close()
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

		setupMongo()
	}

	m.Run()
}

func TestQuery(t *testing.T) {
	Convey("Test Query(ctx *middleware.APIContext, year, make, model string)", t, func() {
		ctx := &middleware.APIContext{
			Session: session,
			DataContext: &customer.DataContext{
				BrandArray: []int{1, 3},
			},
		}

		va, err := Query(ctx, "", "", "")
		So(err, ShouldBeNil)
		So(va.Years, ShouldNotBeEmpty)
		So(len(va.Years), ShouldEqual, 3)

		va, err = Query(ctx, "1997", "", "")
		So(err, ShouldBeNil)
		So(va.Makes, ShouldNotBeEmpty)
		So(len(va.Makes), ShouldEqual, 2)

		va, err = Query(ctx, "1997", "Jeep", "")
		So(err, ShouldBeNil)
		So(va.Models, ShouldNotBeEmpty)
		So(len(va.Models), ShouldEqual, 1)

		va, err = Query(ctx, "1997", "Jeep", "Grand Cherokee")
		So(err, ShouldBeNil)
		So(va.CategoryStyles, ShouldNotBeEmpty)
		So(len(va.CategoryStyles), ShouldEqual, 2)
	})
}

func TestGetYears(t *testing.T) {
	Convey("Test getYears(ctx *middleware.APIContext)", t, func() {
		ctx := &middleware.APIContext{}

		Convey("with no session", func() {
			years, err := getYears(ctx)
			So(err, ShouldNotBeNil)
			So(years, ShouldBeNil)
		})

		ctx.Session = session

		Convey("with no data context", func() {
			years, err := getYears(ctx)
			So(err, ShouldNotBeNil)
			So(years, ShouldBeNil)
		})

		ctx.DataContext = &customer.DataContext{
			BrandArray: []int{1, 3},
		}

		tmp := database.ProductCollectionName
		database.ProductCollectionName = ""

		Convey("with no collection name", func() {
			years, err := getYears(ctx)
			So(err, ShouldNotBeNil)
			So(years, ShouldBeNil)
		})

		database.ProductCollectionName = tmp

		Convey("success", func() {
			years, err := getYears(ctx)
			So(err, ShouldBeNil)
			So(years, ShouldNotBeNil)
		})
	})
}

func TestGetMakes(t *testing.T) {
	Convey("Test getMakes(ctx *middleware.APIContext, year string)", t, func() {
		ctx := &middleware.APIContext{}

		Convey("with no session", func() {
			makes, err := getMakes(ctx, "")
			So(err, ShouldNotBeNil)
			So(makes, ShouldBeNil)
		})

		ctx.Session = session

		Convey("with no data context", func() {
			makes, err := getMakes(ctx, "")
			So(err, ShouldNotBeNil)
			So(makes, ShouldBeNil)
		})

		ctx.DataContext = &customer.DataContext{
			BrandArray: []int{1, 3},
		}

		tmp := database.ProductCollectionName
		database.ProductCollectionName = ""

		Convey("with no collection name", func() {
			makes, err := getMakes(ctx, "")
			So(err, ShouldNotBeNil)
			So(makes, ShouldBeNil)
		})

		database.ProductCollectionName = tmp

		Convey("empty year", func() {
			makes, err := getMakes(ctx, "")
			So(err, ShouldBeNil)
			So(makes, ShouldBeNil)
		})

		Convey("success", func() {
			makes, err := getMakes(ctx, "1997")
			So(err, ShouldBeNil)
			So(makes, ShouldNotBeNil)
		})
	})
}

// func TestReverseLookup(t *testing.T) {
// 	Convey("Test ReverseLookup(*middleware.APIContext, string)", t, func() {
// 		err := database.Init()
// 		So(err, ShouldBeNil)
//
// 		ctx := middleware.APIContext{
// 			Session: database.ProductMongoSession,
// 			DataContext: &customer.DataContext{
// 				BrandArray: []int{1, 3},
// 			},
// 		}
//
// 		Convey("invalid mongo Connection", func() {
// 			res, err := ReverseMongoLookup(emptyCtx, "")
// 			So(err, ShouldNotBeNil)
// 			So(res, ShouldBeNil)
// 		})
//
// 		Convey("valid", func() {
// 			res, err := ReverseMongoLookup(&ctx, "11000")
// 			So(err, ShouldBeNil)
// 			So(res, ShouldHaveSameTypeAs, []VehicleApplication{})
// 		})
// 	})
// }
