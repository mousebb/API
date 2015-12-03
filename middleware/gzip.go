package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
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
func Gzip(h APIHandler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Get(r, apiContext).(*APIContext)
		if ctx == nil {
			ctx = &APIContext{
				Params: httprouter.Params{},
			}
		}
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(rw, r, ctx.Params)
			return
		}
		rw.Header().Set("Vary", "Accept-Encoding")
		rw.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(rw)
		defer gz.Close()
		// gzr := GzipResponseWriter{Writer: gz, ResponseWriter: rw}

		h.ServeHTTP(rw, r, ctx.Params)
		// h.ServeHTTP(gzr, r)
	})
}
