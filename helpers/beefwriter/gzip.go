package beefwriter

import (
	"io"
	"net/http"
)

// NewGzipResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewGzipResponseWriter(w io.Writer, rw http.ResponseWriter) ResponseWriter {
	return &gzipResponseWriter{w, responseWriter{rw, 0, 0, nil, nil}}
}

type gzipResponseWriter struct {
	io.Writer
	responseWriter
}

func (rw *gzipResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.Writer.Write(b)
	rw.size += size
	rw.body = append(rw.body, b...)
	return size, err
}
