package videos_ctlr

import (
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/video"
)

//gets old videos
func DistinctVideos(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.UniqueVideos(ctx)
}

// New videos, literally from the "video_new" table
func Get(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var v video.Video
	var err error

	if v.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = v.Get(ctx)
	return v, err
}

func GetVideoDetails(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var v video.Video
	var err error

	if v.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = v.GetVideoDetails(ctx)
	return v, err
}

func GetAllVideos(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllVideos(ctx)
}

func GetChannel(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var vchan video.Channel
	var err error

	if vchan.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = vchan.Get(ctx)
	return vchan, err
}

func GetAllChannels(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllChannels(ctx)
}

func GetCdn(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var cdn video.CdnFile
	var err error

	if cdn.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = cdn.Get(ctx)
	return cdn, err
}

func GetAllCdns(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllCdnFiles(ctx)
}

func GetVideoType(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var vt video.VideoType
	var err error

	if vt.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = vt.Get(ctx)
	return vt, err
}

func GetAllVideoTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllVideoTypes(ctx)
}

func GetCdnType(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var v video.CdnFileType
	var err error

	if v.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = v.Get(ctx)
	return v, err
}

func GetAllCdnTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllCdnFileTypes(ctx)
}

func GetChannelType(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var v video.ChannelType
	var err error

	if v.ID, err = strconv.Atoi(ctx.Params.ByName("id")); err != nil {
		return nil, err
	}

	err = v.Get(ctx)
	return v, err
}

func GetAllChannelTypes(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	return video.GetAllChannelTypes(ctx)
}
