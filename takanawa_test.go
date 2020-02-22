package takanawa_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa"
)

type mockHandler struct {
	http.Handler
	MockServeHTTP   func(http.ResponseWriter, *http.Request)
	ServeHTTPCalled bool
}

func (m *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ServeHTTPCalled = true
	if m.MockServeHTTP != nil {
		m.MockServeHTTP(w, r)
	}
}

type mockMiddleware struct {
	takanawa.Middleware
	MockApply   func(http.Handler) http.Handler
	ApplyCalled bool
}

func (m *mockMiddleware) Apply(next http.Handler) http.Handler {
	m.ApplyCalled = true
	if m.MockApply != nil {
		return m.MockApply(next)
	} else {
		return nil
	}
}

func TestComposeMiddleware(t *testing.T) {
	tests := []struct {
		mids []takanawa.Middleware
		nil  bool
	}{
		{mids: nil, nil: true},
		{
			mids: []takanawa.Middleware{
				&mockMiddleware{MockApply: func(next http.Handler) http.Handler { return next }},
			},
		},
		{
			mids: []takanawa.Middleware{
				&mockMiddleware{MockApply: func(next http.Handler) http.Handler { return next }},
				&mockMiddleware{MockApply: func(next http.Handler) http.Handler { return next }},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d mids", len(tt.mids)), func(t *testing.T) {
			t.Helper()

			composed := takanawa.ComposeMiddleware(tt.mids...)
			if tt.nil && composed != nil {
				t.Errorf("got %v, want nil", composed)
			}
			if tt.nil {
				return
			}

			var w struct{ http.ResponseWriter }
			var r *http.Request
			handler := &mockHandler{}
			composed.Apply(handler).ServeHTTP(w, r)
			if !handler.ServeHTTPCalled {
				t.Errorf("ServeHTTPCalled: got false, want true")
			}
			for i, mid := range tt.mids {
				if !mid.(*mockMiddleware).ApplyCalled {
					t.Errorf("[%d].ApplyCalled: got false, want true", i)
				}
			}
		})
	}
}
