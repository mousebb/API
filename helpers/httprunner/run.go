package httprunner

import (
	"github.com/curt-labs/GoAPI/controllers/middleware"
	"github.com/curt-labs/GoAPI/helpers/encoding"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Run makes all the requests, prints the summary. It blocks until
// all work is done.
func (b *Runner) Run() {
	b.results = make(chan *result, b.N)
	if b.Output == "" {
		b.bar = newPb(b.N)
	}

	start := time.Now()
	b.run()
	if b.Output == "" {
		b.bar.Finish()
	}

	printReport(b.N, b.results, b.Output, time.Now().Sub(start))
	close(b.results)
}

func (b *Runner) worker(wg *sync.WaitGroup, ch chan *http.Request) {

	for req := range ch {
		// staticOpts := martini.StaticOptions{
		// 	SkipLogging: true,
		// }
		m := martini.New()
		r := martini.NewRouter()

		// martini.StaticOptions.SkipLogging = true
		switch strings.ToUpper(req.Method) {
		case "GET":
			r.Get(b.Req.URL, b.Req.Handler)
		case "POST":
			r.Post(b.Req.URL, b.Req.Handler)
		case "PUT":
			r.Put(b.Req.URL, b.Req.Handler)
		case "PATCH":
			r.Patch(b.Req.URL, b.Req.Handler)
		case "DELETE":
			r.Delete(b.Req.URL, b.Req.Handler)
		case "HEAD":
			r.Head(b.Req.URL, b.Req.Handler)
		default:
			r.Any(b.Req.URL, b.Req.Handler)
		}

		m.Use(render.Renderer())
		m.Use(encoding.MapEncoder)
		m.Use(middleware.Meddler())
		m.MapTo(r, (*martini.Routes)(nil))

		s := time.Now()
		code := 0
		size := int(0)

		response := httptest.NewRecorder()
		m.ServeHTTP(response, req)

		size, _ = strconv.Atoi(response.Header().Get("Content-Length"))
		code = response.Code
		if b.bar != nil {
			b.bar.Increment()
		}
		wg.Done()

		b.results <- &result{
			statusCode:    code,
			duration:      time.Now().Sub(s),
			contentLength: int64(size),
		}
	}
}

func (b *Runner) run() {
	var wg sync.WaitGroup
	wg.Add(b.N)

	var throttle <-chan time.Time
	if b.Qps > 0 {
		throttle = time.Tick(time.Duration(1e6/(b.Qps)) * time.Microsecond)
	}
	jobs := make(chan *http.Request, b.N)
	for i := 0; i < b.C; i++ {
		go func() {
			b.worker(&wg, jobs)
		}()
	}
	for i := 0; i < b.N; i++ {
		if b.Qps > 0 {
			<-throttle
		}
		jobs <- b.Req.Request()
	}
	close(jobs)

	wg.Wait()
}
