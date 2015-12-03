package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger ...
type Logger struct {
	http.Handler
}

func (lh Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	fn := func() {
		log.Printf("Response Time: %s\n", time.Since(start).String())
	}
	defer fn()
}
