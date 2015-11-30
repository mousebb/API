package middleware

import (
	"net/http"

	"time"

	"github.com/curt-labs/API/helpers/nsq"

	"github.com/segmentio/analytics-go"
)

var (
	ExcusedRoutes = []string{"/status", "/customer/auth", "/customer/user", "/new/customer/auth", "/customer/user/register", "/customer/user/resetPassword", "/cartIntegration/priceTypes", "/cartIntegration", "/cache"}
)

func logRequest(r *http.Request, reqTime time.Duration) {

	key := r.Header.Get("key")
	if key == "" {
		vals := r.URL.Query()
		key = vals.Get("key")
	}
	if key == "" {
		key = r.FormValue("key")
	}

	//don't continue if we still don't have a key!
	if key == "" {
		return
	}

	vals := r.URL.Query()
	trkr := analytics.Track{
		Event:      r.URL.String(),
		UserId:     key,
		Properties: make(map[string]interface{}, 0),
	}

	for k, v := range vals {
		trkr.Properties[k] = v
	}

	trkr.Properties["method"] = r.Method
	trkr.Properties["header"] = r.Header
	trkr.Properties["query"] = r.URL.Query().Encode()
	trkr.Properties["referer"] = r.Referer()
	trkr.Properties["userAgent"] = r.UserAgent()
	trkr.Properties["form"] = r.Form
	trkr.Properties["requestTime"] = int64((reqTime.Nanoseconds() * 1000) * 1000)

	go nsq.Push("API_analytics", &trkr)
}
