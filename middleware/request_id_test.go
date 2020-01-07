package middleware_test

import (
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/mock"
	"github.com/kou64yama/takanawa/middleware"
)

func TestRequestID(t *testing.T) {
	t.Run("without "+takanawa.HeaderTakanawaRequestID, func(t *testing.T) {
		t.Helper()

		reqHeader := http.Header{}
		resHeader := http.Header{}
		w := &mock.ResponseWriter{
			MockHeader: func() http.Header { return resHeader },
		}
		r := &http.Request{
			Header: reqHeader,
		}
		middleware.RequestID()(&mock.Handler{}).ServeHTTP(w, r)

		id := resHeader.Get(takanawa.HeaderTakanawaRequestID)
		if len(id) == 0 {
			t.Error("got nil, want not nil")
		}
	})
	t.Run("with "+takanawa.HeaderTakanawaRequestID, func(t *testing.T) {
		t.Helper()

		reqHeader := http.Header{}
		reqHeader.Set(takanawa.HeaderTakanawaRequestID, "e0968622-a057-46e1-95bd-49380695b639")
		resHeader := http.Header{}
		w := &mock.ResponseWriter{
			MockHeader: func() http.Header { return resHeader },
		}
		r := &http.Request{
			Header: reqHeader,
		}
		middleware.RequestID()(&mock.Handler{}).ServeHTTP(w, r)

		id := resHeader.Get(takanawa.HeaderTakanawaRequestID)
		if id != "e0968622-a057-46e1-95bd-49380695b639" {
			t.Errorf("got %q, want %q", id, "e0968622-a057-46e1-95bd-49380695b639")
		}
	})
}
