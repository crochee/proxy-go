// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package middlewares

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"proxy-go/config/dynamic"
	"proxy-go/logger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	xForwardedProto             = "X-Forwarded-Proto"
	xForwardedFor               = "X-Forwarded-For"
	xForwardedHost              = "X-Forwarded-Host"
	xForwardedPort              = "X-Forwarded-Port"
	xForwardedServer            = "X-Forwarded-Server"
	xForwardedURI               = "X-Forwarded-Uri"
	xForwardedMethod            = "X-Forwarded-Method"
	xForwardedTLSClientCert     = "X-Forwarded-Tls-Client-Cert"
	xForwardedTLSClientCertInfo = "X-Forwarded-Tls-Client-Cert-Info"
	xRealIP                     = "X-Real-Ip"
	connection                  = "Connection"
	upgrade                     = "Upgrade"
)

var hostName string

// StatusClientClosedRequest non-standard HTTP status code for client disconnection.
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection.
const StatusClientClosedRequestText = "Client Closed Request"

func BuildProxy(flushInterval time.Duration, medata dynamic.Medata) (http.Handler, error) {
	var (
		item = rand.Int()
		mtx  sync.Mutex
		err  error
	)
	if hostName, err = os.Hostname(); err != nil {
		hostName = "localhost"
	}
	return &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			rewrite(request)
			u := request.URL
			if request.RequestURI != "" {
				parsedURL, err := url.ParseRequestURI(request.RequestURI)
				if err == nil {
					u = parsedURL
				}
			}
			if strings.HasPrefix(u.Path, medata.Path) {
				var pass string
				switch medata.Mode {
				case "random": // 随机
					pass = medata.LocationList[rand.Int()%len(medata.LocationList)].ProxyPass
				case "round_robin": // 轮询
					mtx.Lock()
					pass = medata.LocationList[item%len(medata.LocationList)].ProxyPass
					item++
					mtx.Unlock()
				default:
				}
				if pass != "" {
					parsedURL, err := url.ParseRequestURI(pass)
					if err == nil {
						request.URL.Scheme = parsedURL.Scheme
						request.URL.Opaque = parsedURL.Opaque
						request.URL.User = parsedURL.User
						request.URL.Host = parsedURL.Host

						u.Path = parsedURL.Path
					}
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

func rewrite(req *http.Request) {
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		clientIP = removeIPv6Zone(clientIP)

		if req.Header.Get(xRealIP) == "" {
			req.Header.Set(xRealIP, clientIP)
		}
	}

	xfProto := req.Header.Get(xForwardedProto)
	if xfProto == "" {
		if isWebsocketRequest(req) {
			if req.TLS != nil {
				req.Header.Set(xForwardedProto, "wss")
			} else {
				req.Header.Set(xForwardedProto, "ws")
			}
		} else {
			if req.TLS != nil {
				req.Header.Set(xForwardedProto, "https")
			} else {
				req.Header.Set(xForwardedProto, "http")
			}
		}
	}

	if xfPort := req.Header.Get(xForwardedPort); xfPort == "" {
		req.Header.Set(xForwardedPort, forwardedPort(req))
	}

	if xfHost := req.Header.Get(xForwardedHost); xfHost == "" && req.Host != "" {
		req.Header.Set(xForwardedHost, req.Host)
	}

	if hostName != "" {
		req.Header.Set(xForwardedServer, hostName)
	}
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

func removeIPv6Zone(clientIP string) string {
	return strings.Split(clientIP, "%")[0]
}

// isWebsocketRequest returns whether the specified HTTP request is a websocket handshake request.
func isWebsocketRequest(req *http.Request) bool {
	containsHeader := func(name, value string) bool {
		items := strings.Split(req.Header.Get(name), ",")
		for _, item := range items {
			if value == strings.ToLower(strings.TrimSpace(item)) {
				return true
			}
		}
		return false
	}
	return containsHeader(connection, "upgrade") && containsHeader(upgrade, "websocket")
}

func forwardedPort(req *http.Request) string {
	if req == nil {
		return ""
	}

	if _, port, err := net.SplitHostPort(req.Host); err == nil && port != "" {
		return port
	}

	if req.Header.Get(xForwardedProto) == "https" || req.Header.Get(xForwardedProto) == "wss" {
		return "443"
	}

	if req.TLS != nil {
		return "443"
	}

	return "80"
}
