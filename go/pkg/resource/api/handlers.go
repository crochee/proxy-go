// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/6/5

package api

import (
	"crypto/tls"
	"net/http"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/pkg/logger"
	"github.com/crochee/proxy/pkg/middleware"
	"github.com/crochee/proxy/pkg/middleware/accesslog"
	"github.com/crochee/proxy/pkg/middleware/balance"
	"github.com/crochee/proxy/pkg/middleware/circuitbreaker"
	"github.com/crochee/proxy/pkg/middleware/cros"
	"github.com/crochee/proxy/pkg/middleware/metric"
	"github.com/crochee/proxy/pkg/middleware/ratelimit"
	"github.com/crochee/proxy/pkg/middleware/recovery"
	"github.com/crochee/proxy/pkg/middleware/retry"
	"github.com/crochee/proxy/pkg/middleware/trace"
	"github.com/crochee/proxy/pkg/proxy/httpx"
	"github.com/crochee/proxy/pkg/tlsx"
	"github.com/crochee/proxy/pkg/tracex"
	"github.com/crochee/proxy/version"
)

// Handlers
func Handlers(cfg dynamic.Config) []middleware.Handler {
	var handlers []middleware.Handler
	if cfg.Proxy != nil {
		handlers = append(handlers, proxyHandler(cfg.Proxy))
	}
	if cfg.Middleware == nil {
		return handlers
	}
	if cfg.Middleware.Retry != nil {
		handlers = append(handlers, retry.New(*cfg.Middleware.Retry))
	}
	if cfg.Middleware.CrossDomain {
		handlers = append(handlers, cros.New(cros.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete,
				http.MethodPut, http.MethodPatch, http.MethodHead},
			AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Auth-Token"},
			ExposedHeaders: nil,
			MaxAge:         24 * 60 * 60,
		}))
	}
	if cfg.Middleware.AccessLog != nil {
		handlers = append(handlers, accesslog.New(logger.NewLogger(
			logger.Path(cfg.Middleware.AccessLog.Path),
			logger.Level(cfg.Middleware.AccessLog.Level))))
	}
	if cfg.Middleware.Trace != nil && cfg.Middleware.Trace.Jaeger != nil {
		t, err := tracex.NewTracer(version.ServiceName, 20, cfg.Middleware.Trace.Jaeger)
		if err == nil {
			handlers = append(handlers, trace.NewTraceEntryPoint(t, version.ServiceName))
		} else {
			logger.Warnf("new trace failed.Error:%v", err)
		}
	}
	if cfg.Middleware.Balance != nil {
		handlers = append(handlers, balance.New(*cfg.Middleware.Balance))
	}

	if cfg.Middleware.RateLimit != nil {
		handlers = append(handlers, ratelimit.New(
			ratelimit.Burst(cfg.Middleware.RateLimit.Burst),
			ratelimit.Every(cfg.Middleware.RateLimit.Every),
			ratelimit.Mode(cfg.Middleware.RateLimit.Mode)))
	}
	if cfg.Middleware.CircuitBreaker != nil {
		handlers = append(handlers, circuitbreaker.New(*cfg.Middleware.CircuitBreaker))
	}
	if cfg.Middleware.Recovery {
		handlers = append(handlers, recovery.New())
	}
	handlers = append(handlers, metric.New())
	return handlers
}

func proxyHandler(cfg *dynamic.Proxy) middleware.Handler {
	var proxyOption []httpx.ProxyOption
	if cfg != nil {
		if cfg.Tls != nil {
			tlsConfig, err := tlsx.TlsConfig(tls.RequireAndVerifyClientCert, *cfg.Tls)
			if err == nil {
				proxyOption = append(proxyOption, httpx.TlsConfig(tlsConfig))
			} else {
				logger.Warnf("proxy form https to http.Cause:%v", err)
			}
		}
		if cfg.ProxyLog != nil {
			proxyOption = append(proxyOption, httpx.ProxyLog(logger.NewLogger(logger.Path(cfg.ProxyLog.Path),
				logger.Level(cfg.ProxyLog.Level))))
		}
	}
	return httpx.New(proxyOption...)
}
