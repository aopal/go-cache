package middlewares

import (
	"log"
	"net/http"
)

func WithLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)

		log.Printf("%s %s: %s\n", r.Method, r.URL.String(), w.Header().Get("x-cache-status"))
	}
}
