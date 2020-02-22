package middleware_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/middleware"
)

func TestCors(t *testing.T) {
	tests := []struct {
		host                          string
		option                        *middleware.CorsOption
		accessControlAllowOrigin      string
		accessControlAllowMethods     string
		accessControlAllowHeaders     string
		accessControlExposeHeaders    string
		accessControlAllowCredentials string
	}{
		{
			host:   "localhost",
			option: nil,
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				AllowedOrigins: []string{"*"},
			},
			accessControlAllowOrigin: "*",
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				AllowedOrigins: []string{"localhost", "example.com"},
			},
			accessControlAllowOrigin: "localhost",
		},
		{
			host: "example.com",
			option: &middleware.CorsOption{
				AllowedOrigins: []string{"localhost", "example.com"},
			},
			accessControlAllowOrigin: "example.com",
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				AllowedMethods: []string{"GET", "POST"},
			},
			accessControlAllowMethods: "GET, POST",
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				AllowedHeaders: []string{"Accept", "Content-Type"},
			},
			accessControlAllowHeaders: "Accept, Content-Type",
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				ExposedHeaders: []string{"X-Foo, X-Bar"},
			},
			accessControlExposeHeaders: "X-Foo, X-Bar",
		},
		{
			host: "localhost",
			option: &middleware.CorsOption{
				AllowCredentials: true,
			},
			accessControlAllowCredentials: "true",
		},
	}

	for _, tt := range tests {
		n := "host=" + tt.host
		if tt.option != nil {
			n += fmt.Sprintf(
				",allowedOrigins=%v,allowedMethods=%v,allowedHeaders=%v,exposeHeaders=%v,allowCredential=%v",
				tt.option.AllowedOrigins,
				tt.option.AllowedMethods,
				tt.option.AllowedHeaders,
				tt.option.ExposedHeaders,
				tt.option.AllowCredentials,
			)
		}
		t.Run(n, func(t *testing.T) {
			t.Helper()

			header := http.Header{}
			w := &mock.ResponseWriter{
				MockHeader: func() http.Header { return header },
			}
			r := &http.Request{
				Host: tt.host,
			}
			n := &mock.Handler{}
			middleware.Cors(tt.option).Apply(n).ServeHTTP(w, r)

			accessControlAllowOrigin := header.Get(takanawa.HeaderAccessControlAllowOrigin)
			if accessControlAllowOrigin != tt.accessControlAllowOrigin {
				t.Errorf("got %q, want %q", accessControlAllowOrigin, tt.accessControlAllowOrigin)
			}
			accessControlAllowMethods := header.Get(takanawa.HeaderAccessControlAllowMethods)
			if accessControlAllowMethods != tt.accessControlAllowMethods {
				t.Errorf("got %q, want %q", accessControlAllowMethods, tt.accessControlAllowMethods)
			}
			accessControlExposeHeaders := header.Get(takanawa.HeaderAccessControlExposeHeaders)
			if accessControlExposeHeaders != tt.accessControlExposeHeaders {
				t.Errorf("got %q, want %q", accessControlExposeHeaders, tt.accessControlExposeHeaders)
			}
			accessControlAllowCredentials := header.Get(takanawa.HeaderAccessControlAllowCredentials)
			if accessControlAllowCredentials != tt.accessControlAllowCredentials {
				t.Errorf("got %q, want %q", accessControlAllowCredentials, tt.accessControlAllowCredentials)
			}
		})
	}
}
