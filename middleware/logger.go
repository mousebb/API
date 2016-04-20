package middleware

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	ps "google.golang.org/cloud/pubsub"

	"github.com/curt-labs/API/helpers/beefwriter"
	"github.com/curt-labs/API/helpers/pubsub"
	"github.com/curt-labs/API/models/customer"
)

var (
	topic            = os.Getenv("ANALYTICS_TOPIC")
	analyticsAccount = os.Getenv("ANALYTICS_ACCOUNT")
)

// Header A key-value store for structuring header information. We need
// this since json.Marshal can't handle map[string][]string "net/http Header".
type Header struct {
	Key   string   `bson:"key" json:"key" xml:"key"`
	Value []string `bson:"value" json:"value" xml:"value"`
}

// RequestMetrics Holds data surrounding the incoming request.
type RequestMetrics struct {
	IP          string    `bson:"ip" json:"ip" xml:"ip"`
	ContentType string    `bson:"content_type" json:"content_type" xml:"content_type"`
	Body        []byte    `bson:"body" json:"body" xml:"body"`
	URI         string    `bson:"uri" json:"uri" xml:"uri"`
	Method      string    `bson:"method" json:"method" xml:"method"`
	Headers     []Header  `bson:"headers" json:"headers" xml:"headers"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp" xml:"timestamp"`
}

// ResponseMetrics Holds data surround the outgoing response.
type ResponseMetrics struct {
	ContentType string    `bson:"content_type" json:"content_type" xml:"content_type"`
	StatusCode  int       `bson:"status_code" json:"status_code" xml:"status_code"`
	Headers     []Header  `bson:"headers" json:"headers" xml:"headers"`
	Timestamp   time.Time `bson:"timestamp" json:"timestamp" xml:"timestamp"`
}

// Metrics Holds relevant information that will help report out request
// analytics.
type Metrics struct {
	Application      string          `json:"application" bson:"application"`
	RequestingUser   customer.User   `bson:"user" json:"user" xml:"user"`
	Machine          string          `bson:"machine" json:"machine" xml:"machine"`
	Request          RequestMetrics  `bson:"request_metrics" json:"request_metrics" xml:"request_metrics"`
	Response         ResponseMetrics `bson:"response_metrics" json:"response_metrics" xml:"response_metrics"`
	Latency          int64           `bson:"latency" json:"latency" xml:"latency"`
	Body             []byte          `bson:"body" json:"body" xml:"body"`
	AnalyticsAccount string          `bson:"analytics_account" json:"analytics_account" xml:"analytics_account"`
}

// Log Gathers all information about the request/response and pushes
// the data to Google PubSub.
func Log(w beefwriter.ResponseWriter, r *http.Request, ctx *APIContext) {
	if topic == "" {
		topic = "dev_api_analytics"
	}

	body, _ := ioutil.ReadAll(r.Body)

	var reqHeaders []Header
	for k, v := range r.Header {
		reqHeaders = append(reqHeaders, Header{
			Key:   k,
			Value: v,
		})
	}

	var respHeaders []Header
	for k, v := range w.Header() {
		respHeaders = append(respHeaders, Header{
			Key:   k,
			Value: v,
		})
	}

	reqMetrics := RequestMetrics{
		IP:          r.RemoteAddr,
		ContentType: r.Header.Get("Content-Type"),
		Body:        body,
		URI:         r.URL.String(),
		Method:      r.Method,
		Headers:     reqHeaders,
		Timestamp:   ctx.RequestStart,
	}

	respMetrics := ResponseMetrics{
		ContentType: w.Header().Get("Content-Type"),
		StatusCode:  w.Status(),
		Headers:     respHeaders,
		Timestamp:   time.Now(),
	}

	data := Metrics{
		AnalyticsAccount: analyticsAccount,
		Application:      "apiv3.1",
		Request:          reqMetrics,
		Response:         respMetrics,
		Latency:          time.Since(ctx.RequestStart).Nanoseconds(),
		Body:             w.Body(),
	}

	if ctx != nil && ctx.DataContext != nil {
		data.RequestingUser = ctx.DataContext.User
	}
	data.Machine, _ = os.Hostname()

	var msg ps.Message
	var err error
	msg.Data, err = json.Marshal(&data)
	if err != nil {
		return
	}

	log.Println(topic, analyticsAccount)
	pubsub.PushMessage(topic, &msg)
}
