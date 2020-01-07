package util

import (
	"fmt"
	"net/http"
	"strings"
)

type ResponseSniffer struct {
	http.ResponseWriter
	Writer     http.ResponseWriter
	StatusCode int
	Length     uint64
}

func (w *ResponseSniffer) Header() http.Header {
	return w.Writer.Header()
}

func (w *ResponseSniffer) Write(b []byte) (int, error) {
	size, err := w.Writer.Write(b)
	w.Length += uint64(size)
	return size, err
}

func (w *ResponseSniffer) WriteHeader(statusCode int) {
	w.Writer.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func FormatRequest(req *http.Request) string {
	s := fmt.Sprintln(req.Method, req.RequestURI, req.Proto)
	s += fmt.Sprintln("Host:", req.Host)
	for n, v := range req.Header {
		s += fmt.Sprintf("%s: %s\n", n, strings.Join(v, "; "))
	}
	return s
}

func FormatResponse(resp *http.Response) string {
	s := fmt.Sprintln(resp.Proto, resp.StatusCode)
	for n, v := range resp.Header {
		s += fmt.Sprintf("%s: %s\n", n, strings.Join(v, "; "))
	}
	return s
}
