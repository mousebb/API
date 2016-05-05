package search

import (
	"os"
	"testing"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/customer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDsl(t *testing.T) {
	ctx := &middleware.APIContext{
		DataContext: &customer.DataContext{},
	}

	ip := "127.0.0.1:9200"
	port := "9200"
	if os.Getenv("ELASTIC_HOST") != "" {
		ip = os.Getenv("ELASTIC_HOST")
	}
	if os.Getenv("ELASTIC_PORT") != "" {
		port = os.Getenv("ELASTIC_PORT")
	}
	user := os.Getenv("ELASTIC_USER")
	pass := os.Getenv("ELASTIC_PASS")
	Convey("Testing Search Dsl", t, func() {

		Convey("empty query", func() {
			res, err := Dsl(ctx, "", 0, 0, 0)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})
		Convey("query of `hitch` but bad connections", func() {
			os.Setenv("ELASTIC_HOST", "")
			os.Setenv("ELASTIC_PORT", "")
			os.Setenv("ELASTIC_USER", "")
			os.Setenv("ELASTIC_PASS", "")
			res, err := Dsl(ctx, "hitch", 0, 0, 0)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
			os.Setenv("ELASTIC_HOST", ip)
			os.Setenv("ELASTIC_PORT", port)
			os.Setenv("ELASTIC_USER", user)
			os.Setenv("ELASTIC_PASS", pass)
		})
		Convey("query of `hitch` with no brand", func() {
			res, err := Dsl(ctx, "hitch", 1, 0, 0)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})
		Convey("query of `hitch` with aries", func() {
			_, err := Dsl(ctx, "hitch", 0, 0, 3)
			So(err, ShouldBeNil)
		})
		Convey("query of `hitch` with curt", func() {
			_, err := Dsl(ctx, "hitch", 0, 0, 1)
			So(err, ShouldBeNil)
		})
		Convey("query of `hitch`", func() {
			_, err := Dsl(ctx, "hitch", 0, 0, 3)
			So(err, ShouldBeNil)
		})
		Convey("query of `hitch` with 1 & 3 brand", func() {
			ctx.DataContext.BrandArray = []int{1, 3}
			_, err := Dsl(ctx, "hitch", 0, 1, 0)
			So(err, ShouldBeNil)
		})
		Convey("query of `hitch` with 1 brand", func() {
			_, err := Dsl(ctx, "hitch", 0, 1, 1)
			So(err, ShouldBeNil)
		})

	})
}
