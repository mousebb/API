package vehicle

// import . "github.com/smartystreets/goconvey/convey"

// func TestVehicle(t *testing.T) {
// 	Convey("Testing Vehicle", t, func() {
// 		v := setupDummyVehicle()
// 		Convey("Testing Get Notes", func() {
// 			notes, err := v.GetNotes(13301)
// 			So(len(notes), ShouldBeGreaterThanOrEqualTo, 0)
// 			So(err, ShouldBeNil)
// 		})
// 		Convey("Testing Reverse Lookup", func() {
// 			vehicles, err := ReverseLookup(13301)
// 			So(len(vehicles), ShouldBeGreaterThanOrEqualTo, 0)
// 			So(err, ShouldBeNil)
// 		})
// 	})
// }
//
// func TestMongoVehicle(t *testing.T) {
// 	Convey("Testing Reverse Lookup", t, func() {
// 		vehicles, err := ReverseMongoLookup(2011541)
// 		So(len(vehicles), ShouldBeGreaterThanOrEqualTo, 0)
// 		So(err, ShouldBeNil)
// 	})
// }
//
// func BenchmarkGetNotes(b *testing.B) {
// 	v := setupDummyVehicle()
// 	for i := 0; i < b.N; i++ {
// 		v.GetNotes(13301)
// 	}
// }
//
// func BenchmarkReverseLookup(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		ReverseLookup(13301)
// 	}
// }
//
// func setupDummyVehicle() *Vehicle {
// 	return &Vehicle{
// 		Year:  2010,
// 		Make:  "Chevrolet",
// 		Model: "Silverado 1500",
// 	}
// }
