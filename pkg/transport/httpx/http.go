// Package httpx
package httpx

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/crochee/proxy-go/pkg/logger"
)

type httpServer struct {
	*http.Server
	net.Listener
	ctx context.Context
	option
}

// New new http AppServer
func New(ctx context.Context, host string, handler http.Handler, opts ...Option) (*httpServer, error) {
	srv := &httpServer{
		Server: &http.Server{
			Handler: handler,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		ctx: ctx,
	}
	for _, opt := range opts {
		opt(&srv.option)
	}
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}
	if srv.tlsConfig != nil {
		ln = tls.NewListener(ln, srv.tlsConfig)
	}
	srv.Listener = ln
	if srv.requestLog != nil {
		srv.ConnContext = func(ctx context.Context, c net.Conn) context.Context {
			return logger.Context(ctx, srv.requestLog)
		}
	}
	logger.Infof("listen srv %s", host)
	return srv, nil
}

func (h *httpServer) Name() string {
	return "HTTP(S)"
}

func (h *httpServer) Start() error {
	for _, f := range h.beforeStart {
		if err := f(); err != nil {
			return err
		}
	}
	if err := h.Serve(h.Listener); err != nil {
		return err
	}
	for _, f := range h.afterStart {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}

func (h *httpServer) Stop() error {
	for _, f := range h.beforeStop {
		if err := f(); err != nil {
			return err
		}
	}
	if err := h.Shutdown(h.ctx); err != nil {
		return err
	}
	for _, f := range h.afterStop {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
