package mock

import (
	"net/http"

	"github.com/kou64yama/takanawa"
)

type middleware struct {
	handle takanawa.MiddlewareFunc
}

func (m *middleware) Handle(w http.ResponseWriter, r *http.Request, next takanawa.NextFunc) {
	m.handle(w, r, next)
}

type MiddlewareMock struct {
	Handle        takanawa.MiddlewareFunc
	HandleCalledN int
}

func NewMiddlewareMock() *MiddlewareMock {
	return &MiddlewareMock{
		Handle: func(w http.ResponseWriter, r *http.Request, next takanawa.NextFunc) {
			next()
		},
	}
}

func (mock *MiddlewareMock) Mock() takanawa.Middleware {
	return &middleware{
		handle: func(w http.ResponseWriter, r *http.Request, next takanawa.NextFunc) {
			mock.HandleCalledN++
			mock.Handle(w, r, next)
		},
	}
}
