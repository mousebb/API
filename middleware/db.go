package middleware

import (
	"net/http"

	"github.com/curt-labs/API/helpers/database"
	"github.com/gorilla/context"
)

// Mongo http.HandlerFunc middleware that will initiate
// a MongoDB session and bind to APIContext.
type DB struct {
	http.Handler
}

func (db DB) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := context.Get(r, apiContext).(*APIContext)
	if ctx == nil {
		ctx = &APIContext{}
	}

	sess := database.ProductMongoSession.Copy()
	defer sess.Close()

	ariesSess := database.AriesMongoSession.Copy()
	defer ariesSess.Close()

	ctx.Session = sess
	ctx.AriesSession = ariesSess
	ctx.DB = database.DB

	context.Set(r, apiContext, ctx)
}
