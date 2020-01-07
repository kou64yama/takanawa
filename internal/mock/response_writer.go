package mock

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	MockHeader        func() http.Header
	CalledHeader      [][]interface{}
	MockWrite         func([]byte) (int, error)
	CalledWrite       [][]interface{}
	MockWriteHeader   func(statusCode int)
	CalledWriteHeader [][]interface{}
}

func (w *ResponseWriter) Header() http.Header {
	w.CalledHeader = append(w.CalledHeader, []interface{}{})
	if w.MockHeader != nil {
		return w.MockHeader()
	}
	return nil
}

func (w *ResponseWriter) Write(b []byte) (int, error) {
	w.CalledWrite = append(w.CalledWrite, []interface{}{b})
	if w.MockWrite != nil {
		return w.MockWrite(b)
	}
	return 0, nil
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.CalledWriteHeader = append(w.CalledWriteHeader, []interface{}{statusCode})
	if w.MockWriteHeader != nil {
		w.MockWriteHeader(statusCode)
	}
}
