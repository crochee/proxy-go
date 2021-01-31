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
	"sync"

	"proxy-go/logger"
	"proxy-go/util"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection.
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection.
const StatusClientClosedRequestText = "Client Closed Request"

func NewProxyBuilder(ctx context.Context) http.Handler {
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.RequestURI = "" // Outgoing request should not have RequestURI

			if _, ok := request.Header["User-Agent"]; !ok {
				request.Header.Set("User-Agent", "proxy")
			}

			// Even if the websocket RFC says that headers should be case-insensitive,
			// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
			// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
			// https://tools.ietf.org/html/rfc6455#page-20
			request.Header["Sec-WebSocket-Key"] = request.Header["Sec-Websocket-Key"]
			request.Header["Sec-WebSocket-Extensions"] = request.Header["Sec-Websocket-Extensions"]
			request.Header["Sec-WebSocket-Accept"] = request.Header["Sec-Websocket-Accept"]
			request.Header["Sec-WebSocket-Protocol"] = request.Header["Sec-Websocket-Protocol"]
			request.Header["Sec-WebSocket-Version"] = request.Header["Sec-Websocket-Version"]
			delete(request.Header, "Sec-Websocket-Key")
			delete(request.Header, "Sec-Websocket-Extensions")
			delete(request.Header, "Sec-Websocket-Accept")
			delete(request.Header, "Sec-Websocket-Protocol")
			delete(request.Header, "Sec-Websocket-Version")
		},
		Transport:  http.DefaultTransport,
		BufferPool: newBufferPool(),
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {
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
			log := logger.FromContext(ctx)
			text := statusText(statusCode)
			log.Errorf("%+v '%d %s' caused by: %v", request, statusCode, text, err)
			writer.WriteHeader(statusCode)
			if _, err = writer.Write(util.Bytes(text)); err != nil {
				log.Errorf("Error %v while writing status code", err)
			}
		},
	}
}

func statusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
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
