package video

import (
	"testing"
)

func TestVideo_New(t *testing.T) {
	// var v Video
	// var ch Channel
	// var cdn CdnFile
	// var cdnft CdnFileType
	// var vt VideoType
	// var ct ChannelType
	// var err error
	//
	// //Creates
	// Convey("Testing Video Stuff", t, func() {
	// 	//create chan type
	// 	ct.Name = "test Channel Type"
	// 	err = ct.Create()
	// 	So(err, ShouldBeNil)
	// 	//create chan
	// 	ch.Title = "test title"
	// 	ch.Type = ct
	// 	err = ch.Create()
	// 	So(err, ShouldBeNil)
	// 	//create cdn type
	// 	cdnft.Title = "test cdntype"
	// 	err = cdnft.Create()
	// 	So(err, ShouldBeNil)
	// 	//create cdn
	// 	cdn.ObjectName = "test cdn"
	// 	cdn.Type = cdnft
	// 	err = cdn.Create()
	// 	So(err, ShouldBeNil)
	// 	//create video type
	// 	vt.Name = "test videoType"
	// 	err = vt.Create()
	// 	So(err, ShouldBeNil)
	// 	//create cat
	// 	// cat.Title = "test cat title"
	// 	// err = cat.Create()
	// 	// So(err, ShouldBeNil)
	// 	//create video
	// 	v.Title = "test vid"
	// 	v.Brands = append(v.Brands, brand.Brand{ID: MockedDTX.BrandID}) //matches mocked brand
	// 	// p.ID = 11000                                                    //force part
	//
	// 	v.VideoType = vt
	// 	v.CategoryIds = append(v.CategoryIds, 1)
	// 	v.Channels = append(v.Channels, ch)
	// 	v.Files = append(v.Files, cdn)
	// 	err = v.Create(MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	//update chan type
	// 	ct.Name = "test Channel Type 2"
	// 	err = ct.Update()
	// 	So(err, ShouldBeNil)
	// 	//update chan
	// 	ch.Title = "test title 2"
	// 	err = ch.Update()
	// 	So(err, ShouldBeNil)
	// 	//update cdn type
	// 	cdnft.Title = "test cdntype 2"
	// 	err = cdnft.Update()
	// 	So(err, ShouldBeNil)
	//
	// 	//update video type
	// 	vt.Name = "test videoType 2"
	// 	err = vt.Update()
	// 	So(err, ShouldBeNil)
	// 	//update cdn
	// 	cdn.ObjectName = "test cdn 2"
	// 	err = cdn.Update()
	// 	So(err, ShouldBeNil)
	// 	//update video
	// 	v.Title = "test vid 2"
	//
	// 	v.PartIds = append(v.PartIds, 110001)
	//
	// 	err = v.Update(MockedDTX)
	// 	So(err, ShouldBeNil)
	//
	// 	//get video
	// 	err = v.Get()
	// 	So(err, ShouldBeNil)
	//
	// 	//get details
	// 	err = v.GetVideoDetails()
	// 	So(err, ShouldBeNil)
	// 	So(len(v.PartIds), ShouldBeGreaterThan, 0)
	// 	//get all
	// 	vs, err := GetAllVideos(MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	So(len(vs), ShouldBeGreaterThan, 0)
	// 	//getall part videos
	// 	vs, err = GetPartVideos(110001)
	// 	So(err, ShouldBeNil)
	// 	So(len(vs), ShouldBeGreaterThan, 0)
	// 	//get all channels
	// 	chs, err := v.GetChannels()
	// 	So(err, ShouldBeNil)
	// 	So(len(chs), ShouldBeGreaterThan, 0)
	// 	//get cdns
	// 	cdns, err := v.GetCdnFiles()
	// 	So(err, ShouldBeNil)
	// 	So(len(cdns), ShouldBeGreaterThan, 0)
	// 	//get chan
	// 	err = ch.Get()
	// 	So(err, ShouldBeNil)
	// 	//get cdn
	// 	err = cdn.Get()
	// 	So(err, ShouldBeNil)
	// 	//get cdn type
	// 	err = cdnft.Get()
	// 	So(err, ShouldBeNil)
	// 	//get video type
	// 	err = vt.Get()
	// 	So(err, ShouldBeNil)
	// 	//get chan type
	// 	err = ct.Get()
	// 	So(err, ShouldBeNil)
	// 	//get all chans
	// 	chs, err = GetAllChannels()
	// 	So(err, ShouldBeNil)
	// 	So(len(chs), ShouldBeGreaterThan, 0)
	// 	//get all cdn
	// 	cdns, err = GetAllCdnFiles()
	// 	So(err, ShouldBeNil)
	// 	So(len(cdns), ShouldBeGreaterThan, 0)
	// 	//get all cdn types
	// 	cdnfts, err := GetAllCdnFileTypes()
	// 	So(err, ShouldBeNil)
	// 	So(len(cdnfts), ShouldBeGreaterThan, 0)
	//
	// 	//get all video types
	// 	vts, err := GetAllVideoTypes()
	// 	So(err, ShouldBeNil)
	// 	So(len(vts), ShouldBeGreaterThan, 0)
	// 	//get all file types
	// 	cts, err := GetAllChannelTypes()
	// 	So(err, ShouldBeNil)
	// 	So(len(cts), ShouldBeGreaterThan, 0)
	//
	// 	//delete video
	// 	err = v.Delete(MockedDTX)
	// 	So(err, ShouldBeNil)
	// 	//delete cdn
	// 	err = cdn.Delete()
	// 	So(err, ShouldBeNil)
	// 	//delete chan
	// 	err = ch.Delete()
	// 	So(err, ShouldBeNil)
	// 	//delete chan type
	// 	err = ct.Delete()
	// 	So(err, ShouldBeNil)
	// 	//delete cdn type
	// 	err = cdnft.Delete()
	// 	So(err, ShouldBeNil)
	// 	//delete video type
	// 	err = vt.Delete()
	// 	So(err, ShouldBeNil)
	//
	// 	//delete cat
	// 	// err = cat.Delete(MockedDTX)
	// 	// So(err, ShouldBeNil)
	// 	// //delete video
	// 	// err = v.Delete(MockedDTX)
	// 	// So(err, ShouldBeNil)
	// })

}

func BenchmarkGetAllVideos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllVideos()
	}
}

func BenchmarkGetPartVideos(b *testing.B) {
	// p := 11000
	for i := 0; i < b.N; i++ {
		// GetPartVideos(11000)
	}
}

func BenchmarkGetAllChannels(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllChannels()
	}
}

func BenchmarkGetAllCdnFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllCdnFiles()
	}
}

func BenchmarkGetAllCdnFileTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllCdnFileTypes()
	}
}

func BenchmarkGetAllVideoTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllVideoTypes()
	}
}

func BenchmarkGetAllChannelTypes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// GetAllChannelTypes()
	}
}
