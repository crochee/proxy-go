// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package httpx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
)

func New(cfg *config.TlsConfig) http.Handler {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	if cfg != nil {
		tlsInfo, err := getTls(cfg)
		if err == nil {
			transport.TLSClientConfig = tlsInfo
		}
	}
	return &httputil.ReverseProxy{
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
		Transport:    transport,
		BufferPool:   internal.BufPool,
		ErrorHandler: errHandler,
	}
}

func getTls(cfg *config.TlsConfig) (*tls.Config, error) {
	caPem, err := cfg.Ca.Read()
	if err != nil {
		return nil, err
	}
	var certPem []byte
	if certPem, err = cfg.Cert.Read(); err != nil {
		return nil, err
	}
	var keyPem []byte
	if keyPem, err = cfg.Key.Read(); err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caPem) {
		return nil, errors.New("failed to parse root certificate")
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPem, keyPem); err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:           []tls.Certificate{certificate},
		RootCAs:                caPool,
		CipherSuites:           []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		SessionTicketsDisabled: true,
		MinVersion:             tls.VersionTLS12,
	}, nil
}

func errHandler(writer http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusInternalServerError
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
	log := logger.FromContext(request.Context())
	text := internal.StatusText(statusCode)
	log.Errorf("%+v '%d %s' caused by: %v", request.URL, statusCode, text, err)
	writer.WriteHeader(statusCode)
	if _, err = writer.Write(internal.Bytes(text)); err != nil {
		log.Errorf("Error %v while writing status code", err)
	}
}
