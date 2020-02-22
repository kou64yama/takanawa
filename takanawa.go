// Package takanawa is the HTTP/HTTPS reverse proxy utilities.
package takanawa

import (
	"net/http"
)

// HTTP headers.
const (
	HeaderTakanawaRequestID             = "X-Takanawa-Request-Id"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// Context keys.
var (
	ContextTakanawaRequestID ContextKey
)

// The ContextKey is the type of context key.
type ContextKey string

// A Middleware transforms http.Handler.
type Middleware interface {
	Apply(http.Handler) http.Handler
}

// The MiddlewareFunc type is an adapter to allow the use of ordinary
// functions as takanawa middleware.
type MiddlewareFunc func(http.Handler) http.Handler

// Apply calls f(next)
func (f MiddlewareFunc) Apply(next http.Handler) http.Handler {
	return f(next)
}

// ComposeMiddleware composes mids.
func ComposeMiddleware(mids ...Middleware) Middleware {
	switch len(mids) {
	case 0:
		return nil
	case 1:
		return mids[0]
	default:
		return MiddlewareFunc(func(handler http.Handler) http.Handler {
			return mids[0].Apply(ComposeMiddleware(mids[1:]...).Apply(handler))
		})
	}
}
