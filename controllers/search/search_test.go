package searchCtlr

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/curt-labs/API/helpers/apicontextmock"
	"github.com/curt-labs/API/helpers/error"
	"github.com/curt-labs/API/helpers/httprunner"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	response httptest.ResponseRecorder
)

func TestSearch(t *testing.T) {
	dtx, err := apicontextmock.Mock()
	if err != nil {
		t.Log(err)
	}

	qs := make(url.Values, 0)
	qs.Add("key", dtx.APIKey)

	Convey("Testing Search with empty term", t, func() {
		response = httprunner.Req(Search, "GET", "/search", "/search", &qs)
		So(response.Code, ShouldEqual, 500)
		So(json.Unmarshal(response.Body.Bytes(), &apierror.ApiErr{}), ShouldBeNil)
	})
	Convey("Testing Search with `Hitch`", t, func() {
		response = httprunner.Req(Search, "GET", "/search/:term", "/search/Hitch", &qs)
		So(response.Code, ShouldEqual, 200)
	})
}
