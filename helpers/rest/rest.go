package rest

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// GetPDF Reads a PDF from a provided URL and writes the response to a byte array
func GetPDF(url string) (buf []byte, err error) {

	res, err := http.Get(url)
	if err != nil {
		return
	}
	buf, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	return
}

// IsJsonRequest Detects if the http.Request is sending the body as a JSON object
func IsJsonRequest(r *http.Request) bool {
	ct := strings.ToLower(r.Header.Get("Content-Type"))
	switch ct {
	case "application/json":
		return true
	case "application/x-javascript":
		return true
	case "text/javascript":
		return true
	case "text/x-javascript":
		return true
	case "text/x-json":
		return true
	case "text/json":
		return true
	}

	return false
}
