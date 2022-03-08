package middlewares

import (
	"net/http"
)

func ForceHTML(next http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(rw, r)
	})
}
