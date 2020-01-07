package middleware_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa/middleware"
)

func TestParseReverseProxyOption(t *testing.T) {
	success := []struct {
		in   string
		url  string
		path string
	}{
		{in: "http://localhost:3000", url: "http://localhost:3000", path: ""},
		{in: "/api:http://localhost:8080/v1", url: "http://localhost:8080/v1", path: "/api"},
	}
	failure := []string{
		"/api",
		"http://localhost:3000\t",
	}

	for _, tt := range success {
		t.Run(tt.in, func(t *testing.T) {
			t.Helper()

			u, opt, _ := middleware.ParseReverseProxyOption(tt.in)
			if u.String() != tt.url {
				t.Errorf("got %q, want %q", u, tt.url)
			}
			if opt.Path != tt.path {
				t.Errorf("got %q, want %q", opt.Path, tt.path)
			}
		})
	}
	for _, tt := range failure {
		t.Run(tt, func(t *testing.T) {
			t.Helper()

			_, _, err := middleware.ParseReverseProxyOption(tt)
			if err == nil {
				t.Error("got nil, want not nil")
			}
		})
	}
}

func TestReverseProxy(t *testing.T) {
	tests := []struct {
		proxy           string
		requestURI      string
		proxyRequestURI string
	}{
		{
			proxy:           "http://%s",
			requestURI:      "/",
			proxyRequestURI: "/",
		},
		{
			proxy:           "/api:http://%s/v1",
			requestURI:      "/api/message",
			proxyRequestURI: "/v1/message",
		},
		{
			proxy:           "/api:http://%s/v1",
			requestURI:      "/",
			proxyRequestURI: "",
		},
	}

	for _, tt := range tests {
		n := fmt.Sprintf("proxy=%s,requestURI=%s", tt.proxy, tt.requestURI)
		t.Run(n, func(t *testing.T) {
			t.Helper()

			upLn, _ := net.Listen("tcp", "127.0.0.1:0")
			defer upLn.Close()
			proxyLn, _ := net.Listen("tcp", "127.0.0.1:0")
			defer proxyLn.Close()

			var proxyRequestURI string
			up := &http.Server{
				Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(200)
					proxyRequestURI = req.RequestURI
				}),
			}
			defer up.Close()

			u, opt, _ := middleware.ParseReverseProxyOption(fmt.Sprintf(tt.proxy, upLn.Addr()))
			mid := middleware.ReverseProxy(u, opt)
			proxy := &http.Server{Handler: mid(http.NotFoundHandler())}
			defer proxy.Close()

			go up.Serve(upLn)
			go proxy.Serve(proxyLn)

			http.Get("http://" + proxyLn.Addr().String() + tt.requestURI)

			if proxyRequestURI != tt.proxyRequestURI {
				t.Errorf("got %q, want %q", proxyRequestURI, tt.proxyRequestURI)
			}
		})
	}
}
