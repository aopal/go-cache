package middlewares

import (
	"log"
	"net/http"
	"time"
)

func WithLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(w, r)
		end := time.Now()

		log.Printf("%s %s: %s in %v\n", r.Method, r.URL.String(), w.Header().Get("x-cache-status"), end.Sub(start))
		// fmt.Println(end.Sub(start).Nanoseconds())
	}
}
