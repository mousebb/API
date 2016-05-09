package searchCtlr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/category"
	"github.com/curt-labs/API/models/customer"
	"github.com/curt-labs/API/models/products"
	"github.com/julienschmidt/httprouter"
	"github.com/mattbaird/elastigo/lib"
	"github.com/ory-am/dockertest"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	conn *elastigo.Conn
)

func TestMain(m *testing.M) {
	if os.Getenv("CI") == "" {
		var es dockertest.ContainerID
		var err error
		es, err = dockertest.ConnectToElasticSearch(15, time.Second, func(url string) bool {
			conn = elastigo.NewConn()

			segs := strings.Split(url, ":")
			if len(segs) != 2 {
				log.Fatalf("ElasticSearch connection failed, with address '%s'.", url)
			}

			conn.Domain = segs[0]
			conn.Port = segs[1]

			os.Setenv("ELASTIC_HOST", conn.Domain)
			os.Setenv("ELASTIC_PORT", conn.Port)

			cat := getExampleCategory("1")
			conn.Index("mongo_all", "category", "1", nil, cat)
			conn.Index("mongo_curt", "category", "1", nil, cat)

			part := getExamplePart("1042")
			conn.Index("mongo_all", "part", "1042", nil, part)
			conn.Index("mongo_aries", "part", "1042", nil, part)

			return true
		})

		defer func() {
			os.Clearenv()
			conn.Close()
			es.KillRemove()
		}()

		if err != nil {
			log.Fatal(err)
		}
	} else {
		conn = elastigo.NewConn()

		conn.Domain = os.Getenv("ELASTICSEARCH_PORT_9300_TCP_ADDR")
		conn.Port = os.Getenv("ELASTICSEARCH_PORT_9300_TCP_PORT")

		res, err := http.Get(fmt.Sprintf("http://%s:%s/_cluster/health?pretty=true", conn.Domain, conn.Port))
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(data))

		os.Setenv("ELASTIC_HOST", conn.Domain)
		os.Setenv("ELASTIC_PORT", conn.Port)

		cat := getExampleCategory("1")
		conn.Index("mongo_all", "category", "1", nil, cat)
		conn.Index("mongo_curt", "category", "1", nil, cat)

		part := getExamplePart("1042")
		conn.Index("mongo_all", "part", "1042", nil, part)
		conn.Index("mongo_aries", "part", "1042", nil, part)

		defer func() {
			os.Clearenv()
			conn.Close()
		}()

	}

	m.Run()

}

func TestSearch(t *testing.T) {
	Convey("Testing Search", t, func() {
		ctx := &middleware.APIContext{
			DataContext: &customer.DataContext{
				BrandID: 3,
			},
			Params: httprouter.Params{
				httprouter.Param{
					Key:   "term",
					Value: "1042",
				},
			},
		}

		Convey("with bad brand/part reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.Params[0].Value = "0"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/search/1042", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Search(ctx, rec, req)
			So(err, ShouldNotBeNil)
			So(resp, ShouldBeNil)
		})

		Convey("with good search reference", func() {
			ctx.DataContext.BrandID = 3
			ctx.DataContext.BrandArray = []int{1, 3}
			ctx.Params[0].Value = "1042"
			rec := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "http://localhost:8080/search/1042", nil)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := Search(ctx, rec, req)
			So(err, ShouldBeNil)
			So(resp, ShouldNotBeNil)
		})
	})
}

func getExampleCategory(cat string) category.Category {
	u := fmt.Sprintf("http://api.curtmfg.com/v3/category/%s?key=9300f7bc-2ca6-11e4-8758-42010af0fd79", cat)
	resp, err := http.Get(u)
	if err != nil {
		return category.Category{}
	}
	defer resp.Body.Close()

	var c category.Category
	json.NewDecoder(resp.Body).Decode(&c)

	return c
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
