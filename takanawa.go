package takanawa

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

type NextFunc func()

type MiddlewareFunc func(http.ResponseWriter, *http.Request, NextFunc)

type Middleware interface {
	Handle(http.ResponseWriter, *http.Request, NextFunc)
}

type middleware struct {
	handle MiddlewareFunc
}

func Handle(handle MiddlewareFunc) Middleware {
	return &middleware{handle: handle}
}

func (m *middleware) Handle(w http.ResponseWriter, r *http.Request, next NextFunc) {
	m.handle(w, r, next)
}

func ProxyMiddleware(target *url.URL, overwriteHost bool) Middleware {
	p := httputil.NewSingleHostReverseProxy(target)
	return Handle(func(w http.ResponseWriter, r *http.Request, _ NextFunc) {
		if overwriteHost {
			r.Host = target.Host
		}
		p.ServeHTTP(w, r)
	})
}

func RequestID() Middleware {
	return Handle(func(w http.ResponseWriter, r *http.Request, next NextFunc) {
		id := r.Header.Get(HeaderTakanawaRequestID)
		if len(id) == 0 {
			u, _ := uuid.NewRandom()
			id = u.String()
			r.Header.Set(HeaderTakanawaRequestID, id)
		}
		w.Header().Set(HeaderTakanawaRequestID, id)

		next()
	})
}

type Cors struct {
	AllowOrigin   string
	AllowMethods  []string
	AllowHeaders  []string
	ExposeHeaders []string
}

func CorsMiddleware(cors *Cors) Middleware {
	return Handle(func(w http.ResponseWriter, r *http.Request, next NextFunc) {
		w.Header().Set(HeaderAccessControlAllowOrigin, cors.AllowOrigin)
		w.Header().Set(HeaderAccessControlAllowMethods, strings.Join(cors.AllowMethods, ", "))
		w.Header().Set(HeaderAccessControlAllowHeaders, strings.Join(cors.AllowHeaders, ", "))
		w.Header().Set(HeaderAccessControlExposeHeaders, strings.Join(cors.ExposeHeaders, ", "))
		next()
	})
}

type composer struct {
	serveHTTP http.HandlerFunc
}

func (c *composer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.serveHTTP(w, r)
}

func ComposeMiddlewares(middlewares ...Middleware) http.Handler {
	length := len(middlewares)
	if length == 0 {
		return &composer{
			serveHTTP: func(w http.ResponseWriter, r *http.Request) {},
		}
	}

	head := middlewares[0]
	tail := middlewares[1:len(middlewares)]
	c := ComposeMiddlewares(tail...)

	return &composer{
		serveHTTP: func(w http.ResponseWriter, r *http.Request) {
			head.Handle(w, r, func() {
				c.ServeHTTP(w, r)
			})
		},
	}
}
