package applicationGuide

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/curt-labs/API/middleware"
	"github.com/curt-labs/API/models/applicationGuide"
)

// GetApplicationGuide ...
func GetApplicationGuide(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var err error
	var ag applicationGuide.ApplicationGuide

	ag.ID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return nil, err
	}

	err = ag.Get(ctx.DataContext)
	if err != nil {
		return nil, err
	}

	return ag, nil
}

// GetApplicationGuidesByWebsite ...
func GetApplicationGuidesByWebsite(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	var ag applicationGuide.ApplicationGuide
	var err error
	ag.Website.ID, err = strconv.Atoi(ctx.Params.ByName("id"))
	if err != nil {
		return nil, fmt.Errorf("%s", "failed to parse website identifier")
	}

	return ag.GetBySite(ctx.DataContext)
}

// CreateApplicationGuide ...
// func CreateApplicationGuide(ctx *middleware.APIContext, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
// 	contType := req.Header.Get("Content-Type")
//
// 	var ag applicationGuide.ApplicationGuide
// 	var err error
//
// 	if strings.Contains(contType, "application/json") {
// 		//json
// 		err = json.NewDecoder(req.Body).Decode(&ag)
// 		if err != nil {
// 			return nil, fmt.Errorf("%s", "error decoding request body")
// 		}
// 	} else {
// 		//else, form
// 		ag.Url = req.FormValue("url")
// 		web := req.FormValue("website_id")
// 		ag.FileType = req.FormValue("file_type")
// 		cat := req.FormValue("category_id")
//
// 		if err != nil {
// 			return nil, fmt.Errorf("%s", "error parsing form")
// 		}
// 		if web != "" {
// 			ag.Website.ID, err = strconv.Atoi(web)
// 		}
// 		if cat != "" {
// 			ag.Category.CategoryID, err = strconv.Atoi(cat)
// 		}
// 		if err != nil {
// 			return nil, fmt.Errorf("%s", "error parsing category identifier or website identifer")
// 		}
// 	}
// 	err = ag.Create(ctx.DataContext)
//
// 	return ag, err
// }
//
// // DeleteApplicationGuide ...
// func DeleteApplicationGuide(ctx *middleware.APIContext, req *http.Request) (interface{}, error) {
// 	var err error
// 	var ag applicationGuide.ApplicationGuide
// 	id, err := strconv.Atoi(ctx.Params.ByName("id"))
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", "failed to parse application guides identifier")
// 	}
//
// 	ag.ID = id
// 	err = ag.Delete()
// 	if err != nil {
// 		return nil, fmt.Errorf("%s", "failed to delete application guide")
// 	}
//
// 	return ag, nil
// }
