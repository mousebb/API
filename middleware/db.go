package middleware

import (
	"net/http"

	"github.com/curt-labs/API/helpers/database"
	"github.com/gorilla/context"
)

// DB http.HandlerFunc middleware that will initiate
// a MongoDB session and bind to APIContext.
type DB struct {
	http.Handler
}

func (db DB) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Get(r, apiContext).(*APIContext)
	if ctx == nil {
		ctx = &APIContext{}
	}

	ctx.Session = database.ProductMongoSession.Copy()
	defer ctx.Session.Close()

	ctx.AriesSession = database.AriesMongoSession.Copy()
	defer ctx.AriesSession.Close()

	ctx.DB = database.DB
	defer ctx.DB.Close()

	context.Set(r, apiContext, ctx)
}
