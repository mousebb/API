package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/curt-labs/API/middleware"
)

var (
	start = time.Now()
)

func Index(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://labs.curtmfg.com", http.StatusFound)
	return
}

func Status(ctx *middleware.APIContext, w http.ResponseWriter, r *http.Request) {

	since := time.Since(start)
	secs := since.Seconds()
	run := fmt.Sprintf("running for %g seconds\n", secs)

	w.WriteHeader(200)
	w.Write([]byte(run))
	return
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(""))

	return
}
