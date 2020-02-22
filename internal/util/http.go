package util

import (
	"net/http"
)

// The ResponseSniffer sniff http.ResponseWriter methods.
type ResponseSniffer struct {
	http.ResponseWriter
	Writer     http.ResponseWriter
	StatusCode int
	Length     uint64
}

// Header calls w.Writer.Header()
func (w *ResponseSniffer) Header() http.Header {
	return w.Writer.Header()
}

func (w *ResponseSniffer) Write(b []byte) (int, error) {
	size, err := w.Writer.Write(b)
	w.Length += uint64(size)
	return size, err
}

// WriteHeader calls w.Writer.WriteHeader(statusCode) and saves statusCode.
func (w *ResponseSniffer) WriteHeader(statusCode int) {
	w.Writer.WriteHeader(statusCode)
	w.StatusCode = statusCode
}
