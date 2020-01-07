package takanawa_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/mock"
)

func TestTakanawa(t *testing.T) {
	hello := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("called hello")
			msg := w.Header().Get("X-Test-Message")
			if len(msg) == 0 {
				msg = "hello"
			} else {
				msg += ", hello"
			}
			w.Header().Set("X-Test-Message", msg)
			next.ServeHTTP(w, r)
		})
	}
	world := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Log("called world")
			msg := w.Header().Get("X-Test-Message")
			if len(msg) == 0 {
				msg = "world"
			} else {
				msg += ", world"
			}
			w.Header().Set("X-Test-Message", msg)
			next.ServeHTTP(w, r)
		})
	}

	tests := []struct {
		middleware []takanawa.Middleware
		message    string
	}{
		{middleware: []takanawa.Middleware{hello, world}, message: "hello, world"},
		{middleware: []takanawa.Middleware{hello}, message: "hello"},
		{middleware: []takanawa.Middleware{}, message: ""},
	}

	for _, tt := range tests {
		n := fmt.Sprintf("%d middleware(s)", len(tt.middleware))
		t.Run(n, func(t *testing.T) {
			t.Helper()

			header := http.Header{}
			w := &mock.ResponseWriter{
				MockHeader: func() http.Header { return header },
				MockWrite:  func(b []byte) (int, error) { return len(b), nil },
			}
			r := &http.Request{}

			ta := &takanawa.Takanawa{}
			ta.Middleware(tt.middleware...)
			ta.Handler().ServeHTTP(w, r)

			msg := header.Get("X-Test-Message")
			if msg != tt.message {
				t.Errorf("got %q, want %q", msg, tt.message)
			}
		})
	}
}
