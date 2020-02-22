package middleware

import (
	"net/http"

	"github.com/kou64yama/takanawa"
)

// ChangeOrigin changes the Host request header.
func ChangeOrigin(host string) takanawa.Middleware {
	return takanawa.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r2 := &http.Request{}
			*r2 = *r
			r2.Host = host
			next.ServeHTTP(w, r2)
		})
	})
}
