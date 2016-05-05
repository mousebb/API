package showcase

import (
	"testing"
)

func TestShowcases(t *testing.T) {
	// var s Showcase
	// s.BrandID = MockedDTX.BrandID
	// Convey("Testing Create Showcase", t, func() {
	// 	u, _ := url.Parse("www.test.com")
	// 	i := Image{
	// 		Path: u,
	// 	}
	// 	s.Text = "Test Content"
	// 	s.Images = append(s.Images, i)
	// 	err = s.Create()
	// 	So(err, ShouldBeNil)
	// })
	// Convey("Update", t, func() {
	// 	s.Text = "New Content"
	// 	s.Active = true
	// 	s.Approved = true
	// 	err = s.Update()
	// 	So(err, ShouldBeNil)
	// })
	//
	// Convey("Get showcase", t, func() {
	// 	err = s.Get(MockedDTX)
	// 	So(err, ShouldBeNil)
	// })
	// Convey("GetAll - No paging", t, func() {
	// 	shows, err := GetAllShowcases(0, 1, false, MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	So(len(shows), ShouldBeGreaterThan, 0)
	// })
	//
	// Convey("GetAll - Paged", t, func() {
	// 	shows, err := GetAllShowcases(0, 1, false, MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	So(len(shows), ShouldBeGreaterThan, 0)
	// })
	//
	// Convey("GetAll - randomized", t, func() {
	// 	shows, err := GetAllShowcases(0, 1, true, MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	So(len(shows), ShouldBeGreaterThan, 0)
	//
	// })
	// Convey("Delete", t, func() {
	// 	err = s.Delete()
	// 	So(err, ShouldBeNil)
	// })
}

func BenchmarkGetAllShowcases(b *testing.B) {
	// for i := 0; i < b.N; i++ {
	// 	GetAllShowcases(0, 1, false, MockedDTX)
	// }
}

func BenchmarkGetShowcase(b *testing.B) {
	// show := setupDummyShowcases()
	// for i := 0; i < b.N; i++ {
	// 	b.StopTimer()
	// 	show.Create()
	// 	b.StartTimer()
	// 	show.Get(MockedDTX)
	// 	b.StopTimer()
	// 	show.Delete()
	// }
}

func BenchmarkCreateShowcases(b *testing.B) {
	// show := setupDummyShowcases()
	// for i := 0; i < b.N; i++ {
	// 	b.StartTimer()
	// 	show.Create()
	// 	b.StopTimer()
	// 	show.Delete()
	// }
}

func BenchmarkUpdateShowcases(b *testing.B) {
	// show := setupDummyShowcases()
	// for i := 0; i < b.N; i++ {
	// 	b.StopTimer()
	// 	show.Create()
	// 	b.StartTimer()
	// 	show.Text = "This is a good test."
	// 	show.Update()
	// 	b.StopTimer()
	// 	show.Delete()
	// }
}

func BenchmarkDeleteShowcases(b *testing.B) {
	// show := setupDummyShowcases()
	// for i := 0; i < b.N; i++ {
	// 	b.StopTimer()
	// 	show.Create()
	// 	b.StartTimer()
	// 	show.Delete()
	// }
}

func setupDummyShowcases() *Showcase {
	return &Showcase{
		Rating:    5,
		Title:     "Test Test",
		Text:      "This is a test.",
		Approved:  true,
		Active:    true,
		FirstName: "TESTER",
		LastName:  "TESTER",
		Location:  "Testville, Oklahoma",
		BrandID:   1,
	}
}
