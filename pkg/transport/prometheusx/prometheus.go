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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/metrics"
)

type prometheusAgent struct {
	*http.Server
	ctx context.Context
}

func New(ctx context.Context, host string) *prometheusAgent {
	metrics.DefineMetrics()
	prometheus.MustRegister(metrics.ReqDurHistogramVec, metrics.ReqCodeTotalCounter)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	p := &prometheusAgent{
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

func (p *prometheusAgent) Name() string {
	return "PROMETHEUS_AGENT"
}

func (p *prometheusAgent) Start() error {
	return p.ListenAndServe()
}

func (p *prometheusAgent) Stop() error {
	prometheus.Unregister(metrics.ReqDurHistogramVec)
	prometheus.Unregister(metrics.ReqCodeTotalCounter)
	return p.Shutdown(p.ctx)
}
