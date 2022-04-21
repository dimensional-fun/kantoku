package main

import (
	"net/http"
)

type LogResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w, StatusCode: 200}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.StatusCode = code
	w.ResponseWriter.WriteHeader(code)
}
