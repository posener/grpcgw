package middleware

import (
	"errors"
	"log"
	"net/http"
)

func APILoggerMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logW := newLogResponseWriter(w)
		handler.ServeHTTP(logW, r)
		log.Printf("API called: status=%d verb=%-5s length=%-4d path=%s", logW.Status(), r.Method, logW.ContentLength(), r.URL.Path)
	})
}

// logResponseWriter is a wrapper around ResponseWriter
// that helps fetching the status code and content length of the response
// for logging purposes.
type logResponseWriter struct {
	http.ResponseWriter
	status int
	length int
	error  string
}

func newLogResponseWriter(w http.ResponseWriter) *logResponseWriter {
	return &logResponseWriter{ResponseWriter: w, status: http.StatusOK}
}

func (w *logResponseWriter) Write(data []byte) (int, error) {
	if w.status == http.StatusInternalServerError {
		log.Printf("Inernal server error: %s", data)
		return 0, errors.New("Internal server error")
	}
	length, e := w.ResponseWriter.Write(data)
	w.length += length
	return length, e
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *logResponseWriter) Status() int {
	return w.status
}

func (w *logResponseWriter) ContentLength() int {
	return w.length
}

func (w *logResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *logResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *logResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
