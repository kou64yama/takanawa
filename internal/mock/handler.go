package mock

import "net/http"

type Handler struct {
	http.Handler
	MockServeHTTP   func(w http.ResponseWriter, r *http.Request)
	CalledServeHTTP [][]interface{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.CalledServeHTTP = append(h.CalledServeHTTP, []interface{}{w, r})
	if h.MockServeHTTP != nil {
		h.MockServeHTTP(w, r)
	}
}
