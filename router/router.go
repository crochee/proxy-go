// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"net/http"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware/accesslog"
	"github.com/crochee/proxy-go/middleware/balance"
	"github.com/crochee/proxy-go/middleware/cros"
	"github.com/crochee/proxy-go/middleware/ratelimit"
	"github.com/crochee/proxy-go/middleware/recovery"
	"github.com/crochee/proxy-go/service/proxy/httpx"
)

func Handler(cfg *config.Spec) http.Handler {
	handler := httpx.New(cfg.Proxy)
	// 中间件组合
	if cfg.Middleware != nil {
		if cfg.Middleware.Balance != nil {
			handler = balance.New(cfg.Middleware.Balance, handler)
		}
		if cfg.Middleware.AccessLog != nil {
			handler = accesslog.New(handler, logger.NewLogger(
				logger.Path(cfg.Middleware.AccessLog.Path),
				logger.Level(cfg.Middleware.AccessLog.Level)))
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
	}
	return handler
}
