package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipResponseWriter Wrapper around compress/gzip
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// Gzip Middleware for compressing the response using gzip
func Gzip(h http.Handler, force bool) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && !force {
			h.ServeHTTP(rw, r)
			return
		}
		rw.Header().Set("Vary", "Accept-Encoding")
		rw.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(rw)
		defer gz.Close()
		gzr := GzipResponseWriter{Writer: gz, ResponseWriter: rw}
		h.ServeHTTP(gzr, r)
	})
}
