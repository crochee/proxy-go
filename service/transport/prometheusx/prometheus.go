// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/27

// Package prometheus
package prometheusx

import (
	"context"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type prometheusAgent struct {
	*http.Server
	ctx  context.Context
	host string
	port string
	path string
}

func New(ctx context.Context, host, path string) *prometheusAgent {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())
	return &prometheusAgent{
		Server: &http.Server{
			Addr:    host,
			Handler: mux,
			BaseContext: func(_ net.Listener) context.Context {
				return ctx
			},
		},
		ctx: ctx,
	}
}

func (p *prometheusAgent) Start() error {
	return p.ListenAndServe()
}

func (p *prometheusAgent) Stop() error {
	return p.Shutdown(p.ctx)
}
