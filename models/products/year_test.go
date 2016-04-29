package products

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLookupGetYears(t *testing.T) {
	var l Lookup
	l.Brands = append(l.Brands, 1)
	Convey("Testing GetYears()", t, func() {
		// err := l.GetYears(MockedDTX)
		// So(err, ShouldEqual, nil)
		// So(l.Years, ShouldNotEqual, nil)
		// So(l.Years, ShouldHaveSameTypeAs, []int{})
	})
}
