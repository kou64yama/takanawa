package middleware

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/kou64yama/takanawa"
)

type ReverseProxyOption struct {
	Path          string
	OverwriteHost bool
	RewritePath   func(string) string
	ErrorLog      *log.Logger
}

func ParseReverseProxyOption(s string) (*url.URL, *ReverseProxyOption, error) {
	path := ""
	if strings.HasPrefix(s, "/") {
		sp := strings.SplitN(s, ":", 2)
		if len(sp) != 2 {
			return nil, nil, errors.New("malformed: " + s)
		}
		path = sp[0]
		s = sp[1]
	}

	u, err := url.Parse(s)
	if err != nil {
		return nil, nil, err
	}

	o := &ReverseProxyOption{}
	o.Path = path
	o.RewritePath = func(p string) string { return p[len(path):] }

	return u, o, nil
}

func ReverseProxy(target *url.URL, opt *ReverseProxyOption) takanawa.Middleware {
	if opt == nil {
		opt = &ReverseProxyOption{}
	}

	p := httputil.NewSingleHostReverseProxy(target)
	p.ErrorLog = opt.ErrorLog

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r1 *http.Request) {
			if !strings.HasPrefix(r1.URL.Path, opt.Path) {
				next.ServeHTTP(w, r1)
				return
			}

			r2 := &http.Request{}
			*r2 = *r1

			if opt.OverwriteHost {
				r2.Host = target.Host
			}

			if opt.RewritePath != nil {
				*r2.URL = *r1.URL
				r2.URL.Path = opt.RewritePath(r1.URL.Path)
			}

			id, ok := r1.Context().Value(takanawa.ContextTakanawaRequestID).(string)
			if ok {
				r2.Header.Set(takanawa.HeaderTakanawaRequestID, id)
			}

			p.ServeHTTP(w, r2)
		})
	}
}
