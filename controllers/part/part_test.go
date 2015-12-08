package partCtlr

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/curt-labs/API/models/vehicle"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/products"
	"github.com/curt-labs/API/models/video"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	session *mgo.Session
	db      *sql.DB

	schemas = map[string]string{
		`baseVehicleSchema`: `CREATE TABLE BaseVehicle (ID int(11) NOT NULL AUTO_INCREMENT,AAIABaseVehicleID int(11) DEFAULT NULL,YearID int(11) NOT NULL,MakeID int(11) NOT NULL,ModelID int(11) NOT NULL,PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=25998 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`vcdbVehicleSchema`: `CREATE TABLE vcdb_Vehicle (ID int(11) NOT NULL AUTO_INCREMENT, BaseVehicleID int(11) NOT NULL, SubModelID int(11) DEFAULT NULL, ConfigID int(11) DEFAULT NULL, AppID int(11) DEFAULT NULL, RegionID int(11) NOT NULL DEFAULT '0', PRIMARY KEY (ID)) ENGINE=InnoDB AUTO_INCREMENT=59887 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`vcdbVehiclePartSchema`: `CREATE TABLE vcdb_VehiclePart ( ID int(11) NOT NULL AUTO_INCREMENT,		  VehicleID int(11) NOT NULL,		  PartNumber int(11) NOT NULL,		  PRIMARY KEY (ID)		) ENGINE=InnoDB AUTO_INCREMENT=350523 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`submodelSchema`:       ` CREATE TABLE Submodel (   ID int(11) NOT NULL AUTO_INCREMENT,   AAIASubmodelID int(11) DEFAULT NULL,   SubmodelName varchar(50) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=2037 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`vcdbMakeSchema`:       ` CREATE TABLE vcdb_Make (   ID int(11) NOT NULL AUTO_INCREMENT,   AAIAMakeID int(11) DEFAULT NULL,   MakeName varchar(50) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`vcdbModelSchema`:      ` CREATE TABLE vcdb_Model (   ID int(11) NOT NULL AUTO_INCREMENT,   AAIAModelID int(11) DEFAULT NULL,   ModelName varchar(100) DEFAULT NULL,   VehicleTypeID int(11) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=3922 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`vcdbConfigAttrSchema`: ` CREATE TABLE VehicleConfigAttribute (   ID int(11) NOT NULL AUTO_INCREMENT,   AttributeID int(11) NOT NULL,   VehicleConfigID int(11) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=64582 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`configAttrSchema`:     ` CREATE TABLE ConfigAttribute (   ID int(11) NOT NULL AUTO_INCREMENT,   ConfigAttributeTypeID int(11) NOT NULL,   parentID int(11) NOT NULL,   vcdbID int(11) DEFAULT NULL,   value varchar(255) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=416 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`,
		`configAttrTypeSchema`: ` CREATE TABLE ConfigAttributeType (   ID int(11) NOT NULL AUTO_INCREMENT,   name varchar(100) NOT NULL,   AcesTypeID int(11) DEFAULT NULL,   sort int(11) NOT NULL,   PRIMARY KEY (ID) ) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPACT`}
)

func TestMain(m *testing.M) {

	mongo, err := dockertest.ConnectToMongoDB(15, time.Second, func(url string) bool {
		var err error
		session, err = mgo.Dial(url)
		if err != nil {
			log.Fatalf("MongoDB connection failed, with address '%s'.", url)
		}

		session.SetMode(mgo.Monotonic, true)

		p := getExamplePart("1042")
		p.Identifier = bson.NewObjectId()
		for i := range p.Categories {
			p.Categories[i].Identifier = bson.NewObjectId()
		}

		session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Insert(&p)

		p = getExamplePart("11000")
		p.Identifier = bson.NewObjectId()
		for i := range p.Categories {
			p.Categories[i].Identifier = bson.NewObjectId()
		}

		session.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Insert(&p)

		return session.Ping() == nil
	})
	if err != nil {
		log.Fatal(err)
	}

	mysql, err := dockertest.ConnectToMySQL(15, time.Second*5, func(url string) bool {
		var err error
		db, err = sql.Open("mysql", url)
		if err != nil {
			log.Fatalf("MySQL connection failed, with address '%s'.", url)
		}

		for _, schema := range schemas {
			_, err = db.Exec(schema)
			if err != nil {
				log.Fatal(err)
			}
		}

		return db.Ping() == nil
	})

	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	session.Close()
	db.Close()
	mysql.KillRemove()
	mongo.KillRemove()
}

func TestIdentifiers(t *testing.T) {
	Convey("Testing part.Identifiers", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params:       httprouter.Params{},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		rec := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://localhost:8080/parts/identifiers?brand=1", nil)
		So(err, ShouldBeNil)

		resp, err := Identifiers(ctx, rec, req)
		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, []string{})
	})
}

func TestAll(t *testing.T) {
	Convey("Testing part.All", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params:       httprouter.Params{},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with no page or count paramters", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts", nil)
			So(err, ShouldBeNil)

			resp, err := All(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

		Convey("with page 0 and count size out of bounds", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts?page=0&count=1000", nil)
			So(err, ShouldBeNil)

			resp, err := All(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with page 1 and count size out of bounds", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts?page=1&count=1000", nil)
			So(err, ShouldBeNil)

			resp, err := All(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with page 1 and count size in of bounds", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts?page=1&count=100", nil)
			So(err, ShouldBeNil)

			resp, err := All(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

	})
}

func TestFeatured(t *testing.T) {
	Convey("Testing part.Featured", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params:       httprouter.Params{},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with count paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/featured?count=1", nil)
			So(err, ShouldBeNil)

			resp, err := Featured(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

		Convey("with count paramter out of bounds", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/featured?count=100", nil)
			So(err, ShouldBeNil)

			resp, err := Featured(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with count and bad brand paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/featured?count=1&brand=curt", nil)
			So(err, ShouldBeNil)

			resp, err := Featured(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with count and brand paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/featured?count=1&brand=1", nil)
			So(err, ShouldBeNil)

			resp, err := Featured(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

		Convey("with no page or count paramters", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/featured", nil)
			So(err, ShouldBeNil)

			resp, err := Featured(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

	})
}

func TestLatest(t *testing.T) {
	Convey("Testing part.Latest", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params:       httprouter.Params{},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with count paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/latest?count=1", nil)
			So(err, ShouldBeNil)

			resp, err := Latest(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

		Convey("with count paramter out of bounds", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/latest?count=100", nil)
			So(err, ShouldBeNil)

			resp, err := Latest(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with count and bad brand paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/latest?count=1&brand=curt", nil)
			So(err, ShouldBeNil)

			resp, err := Latest(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with count and brand paramter", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/latest?count=1&brand=1", nil)
			So(err, ShouldBeNil)

			resp, err := Latest(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

		Convey("with no page or count paramters", func() {
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/parts/latest", nil)
			So(err, ShouldBeNil)

			resp, err := Latest(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

	})
}

func TestGet(t *testing.T) {

	Convey("Testing part.Get", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 0
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Get(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Get(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

	})
}

func TestGetRelated(t *testing.T) {

	Convey("Testing part.GetRelated", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/related", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := GetRelated(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Part{})
		})

	})
}

func TestGetWithVehicle(t *testing.T) {

	Convey("Testing part.GetWithVehicle", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/2010/ford/fusion/se", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := GetWithVehicle(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

	})
}

func TestVehicles(t *testing.T) {

	Convey("Testing part.Vehicles", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/0/vehicles", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Vehicles(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/vehicles", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Vehicles(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []vehicle.Vehicle{})
		})

	})
}

func TestImages(t *testing.T) {

	Convey("Testing part.Images", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/images", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Images(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Image{})
		})

	})
}

func TestAttributes(t *testing.T) {

	Convey("Testing part.Attributes", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/attributes", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Attributes(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Attribute{})
		})

	})
}

func TestGetContent(t *testing.T) {

	Convey("Testing part.GetContent", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/content", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := GetContent(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Content{})
		})

	})
}

func TestPackaging(t *testing.T) {

	Convey("Testing part.Packaging", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/0/packaging", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Packaging(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/packaging", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Packaging(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Package{})
		})

	})
}

func TestActiveApprovedReviews(t *testing.T) {

	Convey("Testing part.ActiveApprovedReviews", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/0/reviews", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := ActiveApprovedReviews(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/reviews", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := ActiveApprovedReviews(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []products.Review{})
		})

	})
}

func TestVideos(t *testing.T) {

	Convey("Testing part.Videos", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/0/reviews", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Videos(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/reviews", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Videos(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, []video.Video{})
		})

	})
}

func TestInstallSheet(t *testing.T) {

	Convey("Testing part.InstallSheet", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042.pdf", nil)
			if err != nil {
				t.Fatal(err)
			}

			InstallSheet(ctx, rec, req)
			So(rec.Code, ShouldEqual, 200)
		})

	})
}

func TestCategories(t *testing.T) {

	Convey("Testing part.Categories", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 0
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/categories", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Categories(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/categories", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Categories(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

	})
}

func TestPrices(t *testing.T) {

	Convey("Testing part.Prices", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &middleware.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "part",
					Value: "1042",
				},
			},
			Session:      session,
			AriesSession: session,
			DB:           db,
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 0
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/0/pricing", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Prices(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with proper brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/part/1042/pricing", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Prices(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})

	})
}

func getExamplePart(part string) products.Part {
	u := fmt.Sprintf("http://api.curtmfg.com/v3/part/%s?key=9300f7bc-2ca6-11e4-8758-42010af0fd79", part)
	resp, err := http.Get(u)
	if err != nil {
		return products.Part{}
	}
	defer resp.Body.Close()

	var p products.Part
	json.NewDecoder(resp.Body).Decode(&p)

	return p
}
