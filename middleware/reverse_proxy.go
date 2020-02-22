package middleware

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kou64yama/takanawa"
)

// ReverseProxy returns the middleware that invokes httputil.ReverseProxy.
func ReverseProxy(target *url.URL) takanawa.Middleware {
	p := httputil.NewSingleHostReverseProxy(target)
	p.ErrorLog = log.New(ioutil.Discard, "", 0)
	return takanawa.MiddlewareFunc(func(http.Handler) http.Handler { return p })
}
