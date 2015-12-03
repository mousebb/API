package middleware

import (
	"fmt"
	"net/http"

	"github.com/curt-labs/API/helpers/error"
)

// Recover ...
type Recover struct {
	http.Handler
}

func (rh Recover) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic: %+v\n", err)
			apierror.GenerateError("Internal Server Error", fmt.Errorf("%+v", err), w, r, http.StatusInternalServerError)
			return
		}
	}()
}
