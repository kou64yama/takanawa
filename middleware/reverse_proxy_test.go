package middleware_test

import (
	"net"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/internal/util"
	"github.com/kou64yama/takanawa/middleware"
)

func TestReverseProxy(t *testing.T) {
	t.Helper()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Error(err)
		return
	}
	defer ln.Close()

	srv := http.Server{}
	defer srv.Close()
	closeIdleConnections := make(chan struct{}, 1)
	go func() {
		t.Logf("Listen on %s", ln.Addr())
		if err := srv.Serve(ln); err != http.ErrServerClosed {
			t.Error(err)
		}
		t.Log("Closed")
		close(closeIdleConnections)
	}()

	h := &mock.Handler{}
	w := &mock.ResponseWriter{
		MockHeader: func() http.Header { return http.Header{} },
	}
	r := &http.Request{
		Method: http.MethodGet,
		URL:    util.MustURL("/"),
		Header: http.Header{},
	}
	middleware.ReverseProxy(util.MustURL("http://"+ln.Addr().String())).Apply(h).ServeHTTP(w, r)
	if len(w.CalledWriteHeader) != 1 {
		t.Errorf("called %d, want 1", len(w.CalledWriteHeader))
	}

	srv.Close()
	<-closeIdleConnections
}
