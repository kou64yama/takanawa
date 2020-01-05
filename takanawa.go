package takanawa

import (
	"net/http"
)

const (
	HeaderTakanawaRequestID             = "X-Takanawa-Request-Id"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

var (
	ContextTakanawaRequestID ContextKey
)

type ContextKey string

// Middleware returns a function. That receives http.Handler and
// returns http.Handler.
type Middleware func(http.Handler) http.Handler

// ComposeMiddleware returns a composite Middleware.
func ComposeMiddleware(head Middleware, tail ...Middleware) Middleware {
	if len(tail) == 0 {
		return head
	}

	return func(handler http.Handler) http.Handler {
		return head(ComposeMiddleware(tail[0], tail[1:]...)(handler))
	}
}

type Takanawa struct {
	middleware []Middleware
}

func (t *Takanawa) Middleware(mid ...Middleware) {
	t.middleware = append(t.middleware, mid...)
}

func (t *Takanawa) Handler() http.Handler {
	if len(t.middleware) == 0 {
		return http.NotFoundHandler()
	}
	return ComposeMiddleware(t.middleware[0], t.middleware[1:]...)(http.NotFoundHandler())
}
