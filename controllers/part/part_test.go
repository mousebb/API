package partCtlr

import (
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/curt-labs/API/models/products"
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

		part := products.Part{}
		part.Identifier = bson.NewObjectId()

		db.SetMode(mgo.Monotonic, true)
		db.DB("product_test").C("products").Insert(&part)

		return db.Ping() == nil
	})

	if err != nil {
		log.Fatal(err)
	}

	defer c.KillRemove()
	defer db.Close()

	os.Exit(m.Run())
}

func TestPartNumber(t *testing.T) {

	Convey("Testing part getter", t, func() {
		cnt, err := db.DB("product_test").C("products").Count()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Count: %d\n", cnt)
	})
}
