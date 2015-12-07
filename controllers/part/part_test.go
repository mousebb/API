package partCtlr

import (
	"encoding/json"
	"fmt"
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
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var db *mgo.Session

func TestMain(m *testing.M) {

	c, err := dockertest.ConnectToMongoDB(15, time.Second, func(url string) bool {
		var err error
		db, err = mgo.Dial(url)
		if err != nil {
			log.Fatalf("MongoDB connection failed, with address '%s'.", url)
		}

		p := getExamplePart("1042")
		p.Identifier = bson.NewObjectId()
		for i := range p.Categories {
			p.Categories[i].Identifier = bson.NewObjectId()
		}

		db.SetMode(mgo.Monotonic, true)
		db.DB(database.ProductMongoDatabase).C(database.ProductCollectionName).Insert(&p)

		return db.Ping() == nil
	})

	if err != nil {
		log.Fatal(err)
	}

	m.Run()

	db.Close()
	c.KillRemove()
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
			Session:      db,
			AriesSession: db,
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
			Session:      db,
			AriesSession: db,
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
			Session:      db,
			AriesSession: db,
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
			Session:      db,
			AriesSession: db,
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
			Session:      db,
			AriesSession: db,
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
			Session:      db,
			AriesSession: db,
		}

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
			Session:      db,
			AriesSession: db,
		}

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
			Session:      db,
			AriesSession: db,
		}

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
