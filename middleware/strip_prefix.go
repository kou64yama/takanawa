package middleware

import (
	"net/http"

	"github.com/kou64yama/takanawa"
)

// StripPrefix returns a middleware that invokes http.StripPrefix(prefix, handler).
func StripPrefix(prefix string) takanawa.Middleware {
	return takanawa.MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.StripPrefix(prefix, next)
	})
}
