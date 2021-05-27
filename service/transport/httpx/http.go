// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/6

// Package httpx
package httpx

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
)

type httpServer struct {
	*http.Server
	net.Listener
	ctx context.Context
}

func New(ctx context.Context, medata *config.Medata, handler http.Handler) (*httpServer, error) {
	ln, err := net.Listen("tcp", medata.Host)
	if err != nil {
		return nil, err
	}

	logger.Infof("server with medata:%+v start to run", medata)
	switch strings.ToLower(medata.Scheme) {
	case "http":
	case "https":
		if ln, err = tlsListener(ln, medata); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("scheme is %s", medata.Scheme)
	}
	srv := &httpServer{
		Server: &http.Server{
			Handler: handler,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		Listener: ln,
		ctx:      ctx,
	}
	if medata.RequestLog != nil {
		requestLog := logger.NewLogger(
			logger.Path(medata.RequestLog.Path), logger.Level(medata.RequestLog.Level))
		srv.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
			return logger.Context(ctx, requestLog)
		}
	}
	return srv, nil
}

func (h *httpServer) Start() error {
	return h.Serve(h.Listener)
}

func (h *httpServer) Stop() error {
	return h.Shutdown(h.ctx)
}

func tlsListener(listener net.Listener, medata *config.Medata) (net.Listener, error) {
	if medata.Tls == nil {
		return nil, errors.New("https haven't tls")
	}
	caPEMBlock, err := medata.Tls.Ca.Read()
	if err != nil {
		return nil, err
	}
	var certPEMBlock []byte
	if certPEMBlock, err = medata.Tls.Cert.Read(); err != nil {
		return nil, err
	}
	var keyPEMBlock []byte
	if keyPEMBlock, err = medata.Tls.Key.Read(); err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEMBlock) {
		return nil, errors.New("failed to parse root certificate")
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
		return nil, err
	}

	return tls.NewListener(listener, &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    pool,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		MinVersion:   tls.VersionTLS12,
	}), nil
}
