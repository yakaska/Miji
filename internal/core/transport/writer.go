package transport

import "net/http"

const UninitializedStatusCode = -1

type RespWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *RespWriter {
	return &RespWriter{
		ResponseWriter: w,
		statusCode:     UninitializedStatusCode,
	}
}

func (rw *RespWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *RespWriter) StatusCode() int {
	if rw.statusCode == UninitializedStatusCode {
		panic("uninitialized status code")
	}
	return rw.statusCode
}
