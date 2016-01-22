package vehicle

import (
	"testing"

	"github.com/curt-labs/API/helpers/database"
	"github.com/curt-labs/API/middleware"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGetYears(t *testing.T) {
	Convey("Test GetYears(*middleware.APIContext)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

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

		Convey("valid", func() {
			res, err := GetMakes(&ctx, "2010")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})
		})
	})
}

func TestGetModels(t *testing.T) {
	Convey("Test GetModels(*middleware.APIContext, string)", t, func() {
		err := database.Init()
		So(err, ShouldBeNil)

		ctx := middleware.APIContext{
			Session: database.ProductMongoSession,
		}

		Convey("valid", func() {
			res, err := GetModels(&ctx, "1962", "Ford")
			So(err, ShouldBeNil)
			So(res, ShouldHaveSameTypeAs, []string{})
			t.Log(res)

		})
	})
}
