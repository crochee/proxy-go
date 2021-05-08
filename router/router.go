// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"net/http"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware/accesslog"
	"github.com/crochee/proxy-go/middleware/balance"
	"github.com/crochee/proxy-go/middleware/cros"
	"github.com/crochee/proxy-go/middleware/ratelimit"
	"github.com/crochee/proxy-go/middleware/recovery"
	"github.com/crochee/proxy-go/service/proxy"
)

func Handler(cfg *dynamic.Config) http.Handler {
	handler := proxy.NewProxyBuilder()
	// 中间件组合
	if cfg != nil {
		if len(cfg.Balance) != 0 {
			handler = balance.New(cfg, handler)
		}
		if cfg.AccessLog != nil {
			handler = accesslog.New(handler, logger.NewLogger(
				logger.Path(cfg.AccessLog.Path), logger.Level(cfg.AccessLog.Level)))
		}
		if cfg.RateLimit != nil {
			handler = ratelimit.New(handler,
				ratelimit.Burst(cfg.RateLimit.Burst),
				ratelimit.Every(cfg.RateLimit.Every),
				ratelimit.Mode(cfg.RateLimit.Mode))
		}
		if cfg.Recovery {
			handler = recovery.New(handler)
		}
		if cfg.CrossDomain {
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
