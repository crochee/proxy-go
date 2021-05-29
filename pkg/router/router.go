// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"net/http"

	"github.com/crochee/proxy-go/cmd"
	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/pkg/metrics"
	"github.com/crochee/proxy-go/pkg/middleware/accesslog"
	"github.com/crochee/proxy-go/pkg/middleware/balance"
	"github.com/crochee/proxy-go/pkg/middleware/circuitbreaker"
	"github.com/crochee/proxy-go/pkg/middleware/cros"
	"github.com/crochee/proxy-go/pkg/middleware/metric"
	"github.com/crochee/proxy-go/pkg/middleware/ratelimit"
	"github.com/crochee/proxy-go/pkg/middleware/recovery"
	"github.com/crochee/proxy-go/pkg/middleware/trace"
	"github.com/crochee/proxy-go/pkg/proxy/httpx"
	"github.com/crochee/proxy-go/pkg/tracex"
)

func Handler(cfg *config.Spec) http.Handler {
	handler := httpx.New(cfg.Proxy)
	// 中间件组合
	if cfg.Middleware != nil {
		if cfg.Middleware.CrossDomain {
			handler = cros.New(handler, cros.Options{
				AllowedOrigins: []string{"*"},
				AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete,
					http.MethodPut, http.MethodPatch, http.MethodHead},
				AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Auth-Token"},
				ExposedHeaders: nil,
				MaxAge:         24 * 60 * 60,
			})
		}
		if cfg.Middleware.AccessLog != nil {
			handler = accesslog.New(cmd.ServiceName, logger.NewLogger(
				logger.Path(cfg.Middleware.AccessLog.Path),
				logger.Level(cfg.Middleware.AccessLog.Level)), handler)
		}
		if cfg.Middleware.Trace != nil && cfg.Middleware.Trace.Jaeger != nil {
			t, err := tracex.NewTracer(cmd.ServiceName, 20, cfg.Middleware.Trace.Jaeger)
			if err == nil {
				handler = trace.NewTraceEntryPoint(t, cmd.ServiceName, handler)
			} else {
				logger.Errorf("new trace failed.Error:%v", err)
			}
		}
		if cfg.Middleware.Balance != nil {
			handler = balance.New(*cfg.Middleware.Balance, handler)
		}

		if cfg.Middleware.RateLimit != nil {
			handler = ratelimit.New(handler,
				ratelimit.Burst(cfg.Middleware.RateLimit.Burst),
				ratelimit.Every(cfg.Middleware.RateLimit.Every),
				ratelimit.Mode(cfg.Middleware.RateLimit.Mode))
		}
		if cfg.Middleware.Recovery {
			handler = recovery.New(handler)
		}
		if cfg.Middleware.CircuitBreaker != nil {
			cb, err := circuitbreaker.New(*cfg.Middleware.CircuitBreaker, handler)
			if err != nil {
				logger.Errorf("new circuitbreaker failed.Error:%v", err)
			} else {
				handler = cb
			}
		}
		if cfg.Middleware.Metric {
			handler = metric.New(handler, metrics.ReqDurHistogramVec, metrics.ReqCodeTotalCounter)
		}
	}
	return handler
}
