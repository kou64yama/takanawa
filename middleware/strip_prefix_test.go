package middleware_test

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/internal/util"
	"github.com/kou64yama/takanawa/middleware"
)

func TestStripPrefix(t *testing.T) {
	tests := []struct {
		prefix  string
		request *url.URL
		called  *url.URL
	}{
		{
			prefix:  "/api/v1",
			request: util.MustURL("/api/v1/path/to/resource"),
			called:  util.MustURL("/path/to/resource"),
		},
		{
			prefix:  "/api/v1",
			request: util.MustURL("/path/to/resource"),
		},
	}
	for _, tt := range tests {
		n := fmt.Sprintf("prefix=%q, before=%q, after=%q", tt.prefix, tt.request, tt.called)
		t.Run(n, func(t *testing.T) {
			t.Helper()

			h := &mock.Handler{}
			w := &mock.ResponseWriter{
				MockHeader: func() http.Header { return http.Header{} },
			}
			r := &http.Request{
				Method: http.MethodGet,
				URL:    tt.request,
				Header: http.Header{},
			}
			middleware.StripPrefix(tt.prefix).Apply(h).ServeHTTP(w, r)

			if tt.called != nil {
				if len(h.CalledServeHTTP) != 1 {
					t.Errorf("called %d, want 1", len(h.CalledServeHTTP))
				}
				req := h.CalledServeHTTP[0][1].(*http.Request)
				if req.URL.String() != tt.called.String() {
					t.Errorf("got %s, want %s", req.URL, tt.called)
				}
			} else {
				if len(h.CalledServeHTTP) > 0 {
					t.Errorf("called %d, want 0", len(h.CalledServeHTTP))
				}
			}
		})
	}
}
