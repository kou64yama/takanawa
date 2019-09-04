package takanawa_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kou64yama/takanawa"
	"github.com/kou64yama/takanawa/internal/assert"
	"github.com/kou64yama/takanawa/internal/mock"
)

func TestUpstream(t *testing.T) {

}

func TestProxy(t *testing.T) {

}

func TestRequestID(t *testing.T) {
	tests := []string{
		"1d402cb9-1149-4d5f-88d5-425fdbb2922c",
		"",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			middleware := takanawa.RequestID()

			next := mock.NewNextMock()
			w := mock.NewResponseWriterMock()
			h := http.Header{}
			h.Set(takanawa.HeaderTakanawaRequestID, tt)
			r := &http.Request{
				Header: h,
			}

			middleware.Handle(w.Mock(), r, next.Mock())

			ass := assert.NewAssertions(t)
			ass.AssertEquals(next.CalledN, 1)

			reqID := r.Header.Get(takanawa.HeaderTakanawaRequestID)
			resID := w.Header.Get(takanawa.HeaderTakanawaRequestID)
			ass.AssertEquals(reqID, resID)
			ass.AssertTrue(len(reqID) > 0)
			if len(tt) > 0 {
				ass.AssertEquals(reqID, tt)
			}
		})
	}
}

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		cors          *takanawa.Cors
		allowOrigin   string
		allowMethods  string
		allowHeaders  string
		exposeHeaders string
	}{
		{
			cors: &takanawa.Cors{
				AllowOrigin:   "*",
				AllowMethods:  []string{"GET", "POST"},
				AllowHeaders:  []string{"Content-Type", "Content-Length"},
				ExposeHeaders: []string{"X-Meta"},
			},
			allowOrigin:   "*",
			allowMethods:  "GET, POST",
			allowHeaders:  "Content-Type, Content-Length",
			exposeHeaders: "X-Meta",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.cors), func(t *testing.T) {
			m := takanawa.CorsMiddleware(tt.cors)

			next := mock.NewNextMock()
			w := mock.NewResponseWriterMock()

			m.Handle(w.Mock(), nil, next.Mock())

			ass := assert.NewAssertions(t)
			ass.AssertEquals(next.CalledN, 1)
			ass.AssertEquals(w.Header.Get(takanawa.HeaderAccessControlAllowOrigin), tt.allowOrigin)
			ass.AssertEquals(w.Header.Get(takanawa.HeaderAccessControlAllowMethods), tt.allowMethods)
			ass.AssertEquals(w.Header.Get(takanawa.HeaderAccessControlAllowHeaders), tt.allowHeaders)
			ass.AssertEquals(w.Header.Get(takanawa.HeaderAccessControlExposeHeaders), tt.exposeHeaders)
		})
	}
}

func TestComposeMiddlewares(t *testing.T) {
	tests := [][]*mock.MiddlewareMock{
		{},
		{mock.NewMiddlewareMock()},
		{mock.NewMiddlewareMock(), mock.NewMiddlewareMock()},
		{mock.NewMiddlewareMock(), mock.NewMiddlewareMock(), mock.NewMiddlewareMock()},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d middleware(s)", len(tt)), func(t *testing.T) {
			middlewares := make([]takanawa.Middleware, len(tt))
			for i, m := range tt {
				middlewares[i] = m.Mock()
			}
			composed := takanawa.ComposeMiddlewares(middlewares...)

			w := mock.NewResponseWriterMock()
			r := &http.Request{}
			composed.ServeHTTP(w.Mock(), r)

			ass := assert.NewAssertions(t)
			for _, m := range tt {
				ass.AssertEquals(m.HandleCalledN, 1)
			}
		})
	}
}
