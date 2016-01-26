package vehicle

import (
	"testing"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	emptyCtx = middleware.APIContext{}
)

func TestReverseLookup(t *testing.T) {
	Convey("Test ReverseLookup(*middleware.APIContext, string)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("invalid mongo Connection", func() {
			res, err := ReverseMongoLookup(&emptyCtx, "")
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})

		Convey("valid", func() {
			res, err := ReverseMongoLookup(&ctx, "11000")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []VehicleApplication{})
		})
	})
}

func TestGetYears(t *testing.T) {
	Convey("Test GetYears(*middleware.APIContext)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("invalid mongo Connection", func() {
			res, err := GetYears(&emptyCtx)
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})

		Convey("valid", func() {
			res, err := GetYears(&ctx)
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetMakes(t *testing.T) {
	Convey("Test GetMakes(*middleware.APIContext, string)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("invalid mongo Connection", func() {
			res, err := GetMakes(&emptyCtx, "")
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})

		Convey("valid", func() {
			res, err := GetMakes(&ctx, "1962")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetModels(t *testing.T) {
	Convey("Test GetModels(*middleware.APIContext, string, string)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("invalid mongo Connection", func() {
			res, err := GetModels(&emptyCtx, "", "")
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})

		Convey("valid", func() {
			res, err := GetModels(&ctx, "1962", "Chevrolet")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})

		})
	})
}

func TestGetStyles(t *testing.T) {
	Convey("Test GetStyles(*middleware.APIContext, string, string, string)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("invalid mongo Connection", func() {
			res, err := GetStyles(&emptyCtx, "", "", "")
			So(err, ShouldNotBeNil)
			So(res, ShouldBeNil)
		})

		Convey("valid", func() {
			res, err := GetStyles(&ctx, "1962", "Chevrolet", "All Full Size Pickups")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})

		})
	})
}

func BenchmarkGetStylesFromCollections(b *testing.B) {
	err := database.Init()
	if err != nil {
		b.Fatal(err.Error())
	}

	ctx := middleware.APIContext{
		Session: database.ProductMongoSession,
	}

	var res []string
	for i := 0; i < b.N; i++ {
		res, err = GetStylesFromCollections(&ctx, "2010", "Ford", "F-150")
		if err != nil {
			b.Error(err.Error())
		} else if len(res) == 0 {
			b.Error("result was empty")
		}
	}
	b.Log(res)
}

func BenchmarkGetStyles(b *testing.B) {
	err := database.Init()
	if err != nil {
		b.Fatal(err.Error())
	}

	ctx := middleware.APIContext{
		Session: database.ProductMongoSession,
	}

	var res []string
	for i := 0; i < b.N; i++ {
		res, err = GetStyles(&ctx, "2010", "ford", "f-150")
		if err != nil {
			b.Error(err.Error())
		} else if len(res) == 0 {
			b.Error("result was empty")
		}
	}
	b.Log(res)
}
