package mock

import "net/http"

type responseWriter struct {
	header      func() http.Header
	writeHeader func(int)
	write       func([]byte) (int, error)
}

func (w *responseWriter) Header() http.Header {
	return w.header()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.writeHeader(statusCode)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.write(b)
}

type ResponseWriterMock struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func NewResponseWriterMock() *ResponseWriterMock {
	return &ResponseWriterMock{
		Header: http.Header{},
	}
}

func (mock *ResponseWriterMock) Mock() http.ResponseWriter {
	return &responseWriter{
		header: func() http.Header {
			return mock.Header
		},
		writeHeader: func(statusCode int) {
			mock.StatusCode = statusCode
		},
		write: func(b []byte) (int, error) {
			mock.Body = append(mock.Body, b...)
			return len(b), nil
		},
	}
}
