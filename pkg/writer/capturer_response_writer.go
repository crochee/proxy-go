// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package writer

import "net/http"

type Capture interface {
	http.ResponseWriter
	Size() int64
	Status() int
}

// captureResponseWriter is a wrapper of type http.ResponseWriter
// that tracks request status and size.
type captureResponseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int64
}

func (c *captureResponseWriter) Header() http.Header {
	return c.rw.Header()
}

func (c *captureResponseWriter) Write(bytes []byte) (int, error) {
	if c.status == 0 {
		c.status = http.StatusOK
	}
	size, err := c.rw.Write(bytes)
	c.size += int64(size)
	return size, err
}

func (c *captureResponseWriter) WriteHeader(statusCode int) {
	c.rw.WriteHeader(statusCode)
	c.status = statusCode
}

func (c *captureResponseWriter) Size() int64 {
	return c.size
}

func (c *captureResponseWriter) Status() int {
	return c.status
}

type captureResponseWriterWithCloseNotify struct {
	*captureResponseWriter
}

func NewCaptureResponseWriter(rw http.ResponseWriter) Capture {
	capt := &captureResponseWriter{rw: rw}
	if _, ok := rw.(http.CloseNotifier); !ok { // nolint:staticcheck
		return capt
	}
	return &captureResponseWriterWithCloseNotify{capt}
}

func (c *captureResponseWriterWithCloseNotify) CloseNotify() <-chan bool {
	return c.rw.(http.CloseNotifier).CloseNotify() // nolint:staticcheck
}
