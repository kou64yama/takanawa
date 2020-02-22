package middleware_test

import (
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/middleware"
)

func TestChangeOrigin(t *testing.T) {
	t.Helper()

	h := &mock.Handler{}
	w := &mock.ResponseWriter{}
	r := &http.Request{Host: "localhost"}
	middleware.ChangeOrigin("example.com").Apply(h).ServeHTTP(w, r)

	if len(h.CalledServeHTTP) == 0 {
		t.Error("handler not called, want called")
		return
	}
	req := h.CalledServeHTTP[0][1].(*http.Request)
	if req.Host != "example.com" {
		t.Errorf("got %s, want %s", req.Host, "example.com")
	}
}
