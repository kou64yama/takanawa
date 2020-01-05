package middleware

import (
	"net/http"
	"strings"

	"github.com/kou64yama/takanawa"
)

type CorsOption struct {
	takanawa.Middleware
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
}

func Cors(opt *CorsOption) takanawa.Middleware {
	if opt == nil {
		opt = &CorsOption{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			for _, v := range opt.AllowedOrigins {
				if v == "*" || v == r.Host {
					h.Set(takanawa.HeaderAccessControlAllowOrigin, v)
					break
				}
			}
			if len(opt.AllowedMethods) > 0 {
				h.Set(takanawa.HeaderAccessControlAllowMethods, strings.Join(opt.AllowedMethods, ", "))
			}
			if len(opt.AllowedHeaders) > 0 {
				h.Set(takanawa.HeaderAccessControlAllowHeaders, strings.Join(opt.AllowedHeaders, ", "))
			}
			if len(opt.ExposedHeaders) > 0 {
				h.Set(takanawa.HeaderAccessControlExposeHeaders, strings.Join(opt.ExposedHeaders, ", "))
			}
			if opt.AllowCredentials {
				h.Set(takanawa.HeaderAccessControlAllowCredentials, "true")
			}

			next.ServeHTTP(w, r)
		})
	}
}
