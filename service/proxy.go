// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package service

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"proxy-go/logger"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection.
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection.
const StatusClientClosedRequestText = "Client Closed Request"

func NewProxyBuilder(flushInterval time.Duration) (http.Handler, error) {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			u := request.URL
			if request.RequestURI != "" {
				parsedURL, err := url.ParseRequestURI(request.RequestURI)
				if err == nil {
					u = parsedURL
				}
			}

			request.URL.Path = u.Path
			request.URL.RawPath = u.RawPath
			request.URL.RawQuery = u.RawQuery
			request.RequestURI = "" // Outgoing request should not have RequestURI
			request.Proto = "HTTP/1.1"
			request.ProtoMajor = 1
			request.ProtoMinor = 1

			if _, ok := request.Header["User-Agent"]; !ok {
				request.Header.Set("User-Agent", "")
			}
		},
		Transport:     http.DefaultTransport,
		FlushInterval: flushInterval,
		BufferPool:    newBufferPool(),
		ErrorHandler:  ErrorHandler,
	}, nil
}

func statusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}

func ErrorHandler(w http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, io.EOF):
		statusCode = http.StatusBadGateway
	case errors.Is(err, context.Canceled):
		statusCode = StatusClientClosedRequest
	default:
		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				statusCode = http.StatusGatewayTimeout
			} else {
				statusCode = http.StatusBadGateway
			}
		}
	}

	logger.Errorf("url:%+v '%d %s' caused by: %v",
		request,
		statusCode, statusText(statusCode), err)
	w.WriteHeader(statusCode)
	if _, err = w.Write([]byte(statusText(statusCode))); err != nil {
		logger.Errorf("Error while writing status code", err)
	}
}

const bufferPoolSize = 32 * 1024

func newBufferPool() *bufferPool {
	return &bufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, bufferPoolSize)
			},
		},
	}
}

type bufferPool struct {
	pool sync.Pool
}

func (b *bufferPool) Get() []byte {
	return b.pool.Get().([]byte)
}

func (b *bufferPool) Put(bytes []byte) {
	b.pool.Put(bytes)
}
