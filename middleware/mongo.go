package middleware

import (
	"net/http"

	"github.com/curt-labs/API/helpers/database"
	"github.com/gorilla/context"
)

// Mongo http.HandlerFunc middleware that will initiate
// a MongoDB session and bind to APIContext.
type Mongo struct {
	http.Handler
}

func (mh Mongo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Get(r, apiContext).(*APIContext)
	if ctx == nil {
		ctx = &APIContext{}
	}

	sess := database.ProductMongoSession.Copy()
	defer sess.Close()

	ctx.Session = sess

	context.Set(r, apiContext, ctx)
}
