package middleware

// Chain Loops over the declared handlers
// func Chain(name string, method string, handlers ...http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
//
// 		for _, h := range handlers {
// 			h.ServeHTTP(w, r)
// 			// if err != nil {
// 			// 	apierror.GenerateError(err.Error(), err, w, r, http.StatusInternalServerError)
// 			// 	return
// 			// }
// 		}
//
// 		end := time.Now()
// 		fmt.Printf("%s\t%s_%s\t%s\n", name, method, r.URL.String(), end.Sub(start))
// 	})
// }
