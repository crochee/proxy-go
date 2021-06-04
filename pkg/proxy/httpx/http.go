package httpx

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/pkg/logger"
)

type proxy struct {
	reverseProxy *httputil.ReverseProxy
}

func (p *proxy) Name() string {
	return "HTTP(S)_PROXY"
}

func (p *proxy) Level() int {
	return 0
}

func (p *proxy) Next(handler http.Handler) {
}

func (p *proxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	p.reverseProxy.ServeHTTP(rw, req)
}

func New(opts ...ProxyOption) *proxy {
	var o option
	for _, opt := range opts {
		opt(&o)
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig:       o.tlsConfig,
		TLSHandshakeTimeout:   10 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}
	return &proxy{reverseProxy: &httputil.ReverseProxy{
		Director: func(request *http.Request) {
			request.RequestURI = "" // Outgoing request should not have RequestURI

			if _, ok := request.Header["User-Agent"]; !ok {
				request.Header.Set("User-Agent", "proxy")
			}
			// Even if the websocket RFC says that headers should be case-insensitive,
			// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
			// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
			delete(request.Header, "Sec-Websocket-Key")
			delete(request.Header, "Sec-Websocket-Extensions")
			delete(request.Header, "Sec-Websocket-Accept")
			delete(request.Header, "Sec-Websocket-Protocol")
			delete(request.Header, "Sec-Websocket-Version")
		},
		Transport:  transport,
		BufferPool: internal.BufPool,
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			errHandler(o.proxyLog, rw, req, err)
		},
	}}
}

func errHandler(log logger.Builder, rw http.ResponseWriter, req *http.Request, err error) {
	var statusCode int
	switch {
	case errors.Is(err, io.EOF):
		statusCode = http.StatusBadGateway
	case errors.Is(err, context.Canceled):
		statusCode = internal.StatusClientClosedRequest
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
	text := internal.StatusText(statusCode)
	if log != nil {
		log.Errorf("[PROXY] %+v '%d %s' caused by: %v", req.URL, statusCode, text, err)
	}
	rw.WriteHeader(statusCode)
	if _, err = rw.Write(internal.Bytes(text)); err != nil && log != nil {
		log.Errorf("[PROXY] Error %v while writing status code", err)
	}
}

type ProxyOption func(*option)

type option struct {
	tlsConfig *tls.Config
	proxyLog  logger.Builder
}

// TlsConfig
func TlsConfig(cfg *tls.Config) ProxyOption {
	return func(o *option) { o.tlsConfig = cfg }
}

// ProxyLog
func ProxyLog(log logger.Builder) ProxyOption {
	return func(o *option) { o.proxyLog = log }
}
