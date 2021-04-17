package main

import (
	"log"
	"net/http"
	"time"
)

// https://www.codemio.com/2019/01/advanced-golang-tutorial-http-middleware.html

// middleware provides a convenient mechanism for filtering HTTP requests
// entering the application. It returns a new handler which may perform various
// operations and should finish by calling the next HTTP handler.
type middleware func(next http.HandlerFunc) http.HandlerFunc

// Logging middleware
func withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s - %s - %dms", r.RemoteAddr, r.RequestURI,
				time.Since(start).Milliseconds())
		}()

		next.ServeHTTP(w, r)
	}
}

// nestedMiddleware provides syntactic sugar to create a new middleware
// which will be the result of chaining the ones received as parameters.
func nestedMiddleware(mw ...middleware) middleware {
	return func(final http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			last := final
			for i := len(mw) - 1; i >= 0; i-- {
				last = mw[i](last)
			}
			last(w, r)
		}
	}
}

// Código propio

// allMiddlewares es el resultado de ejecutad nestedMiddleware una vez.
//Contiene todos los middleware a ejecutar antes de handler recibido en HandleFunc
var allMiddlewares middleware

// HandleFunc Esta función ejecutará todos los middleware antes del handler recibido
func HandleFunc(pattern string, handler http.HandlerFunc) {
	if allMiddlewares == nil {
		allMiddlewares = nestedMiddleware(withLogging)
	}
	http.Handle(pattern, allMiddlewares(handler))
}
