// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/30

package pprofx

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/crochee/proxy-go/pkg/logger"
)

type pprofAgent struct {
	*http.Server
	ctx context.Context
}

// New create pprof server
func New(ctx context.Context, host string) *pprofAgent {
	mux := http.NewServeMux()
	mux.HandleFunc("/", pprof.Index)
	mux.HandleFunc("/cmdline", pprof.Cmdline)
	mux.HandleFunc("/profile", pprof.Profile)
	mux.HandleFunc("/symbol", pprof.Symbol)
	mux.HandleFunc("/trace", pprof.Trace)
	mux.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
	mux.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
	mux.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
	mux.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
	mux.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
	mux.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
	p := &pprofAgent{
		Server: &http.Server{
			Addr:    host,
			Handler: mux,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		ctx: ctx,
	}
	logger.Infof("listen %s,host %s", p.Name(), host)
	return p
}
func (p *pprofAgent) Name() string {
	return "PPROF"
}

func (p *pprofAgent) Start() error {
	return p.ListenAndServe()
}

func (p *pprofAgent) Stop() error {
	return p.Shutdown(p.ctx)
}
