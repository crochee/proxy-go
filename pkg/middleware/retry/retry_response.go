package retry

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type responseWriter interface {
	http.ResponseWriter
	http.Flusher
	ShouldRetry() bool
	DisableRetries()
}

func newResponseWriter(rw http.ResponseWriter, shouldRetry bool) responseWriter {
	retryRw := &responseWriterWithoutCloseNotify{
		responseWriter: rw,
		headers:        make(http.Header),
		shouldRetry:    shouldRetry,
	}
	if _, ok := rw.(http.CloseNotifier); ok { // nolint:staticcheck
		return &responseWriterWithCloseNotify{
			responseWriterWithoutCloseNotify: retryRw,
		}
	}
	return retryRw
}

type responseWriterWithoutCloseNotify struct {
	responseWriter http.ResponseWriter
	headers        http.Header
	shouldRetry    bool
	written        bool
}

func (r *responseWriterWithoutCloseNotify) ShouldRetry() bool {
	return r.shouldRetry
}

func (r *responseWriterWithoutCloseNotify) DisableRetries() {
	r.shouldRetry = false
}

func (r *responseWriterWithoutCloseNotify) Header() http.Header {
	if r.written {
		return r.responseWriter.Header()
	}
	return r.headers
}

func (r *responseWriterWithoutCloseNotify) Write(buf []byte) (int, error) {
	if r.ShouldRetry() {
		return len(buf), nil
	}
	return r.responseWriter.Write(buf)
}

func (r *responseWriterWithoutCloseNotify) WriteHeader(code int) {
	if r.ShouldRetry() {
		// 不重试的机制
		if code == http.StatusServiceUnavailable || code/500 != 1 {
			r.DisableRetries()
		}
	}

	if r.ShouldRetry() {
		return
	}

	// retry header copy or add key-value which fix header
	headers := r.responseWriter.Header()
	for header, value := range r.headers {
		headers[header] = value
	}

	r.responseWriter.WriteHeader(code)
	r.written = true
}

func (r *responseWriterWithoutCloseNotify) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := r.responseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("%T is not a http.Hijacker", r.responseWriter)
	}
	return hijacker.Hijack()
}

func (r *responseWriterWithoutCloseNotify) Flush() {
	if flusher, ok := r.responseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

type responseWriterWithCloseNotify struct {
	*responseWriterWithoutCloseNotify
}

func (r *responseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return r.responseWriter.(http.CloseNotifier).CloseNotify() // nolint:staticcheck
}
