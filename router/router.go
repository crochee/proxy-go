// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"github.com/crochee/proxy-go/middleware/ratelimit"
	"net/http"
	"regexp"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/middleware/balance"
	"github.com/crochee/proxy-go/middleware/cros"
	"github.com/crochee/proxy-go/middleware/logger"
	"github.com/crochee/proxy-go/middleware/mix"
	"github.com/crochee/proxy-go/middleware/recovery"
	"github.com/crochee/proxy-go/service"
)

func ChainBuilder() (http.Handler, error) {
	proxy := service.NewProxyBuilder()

	balanceProxy, ok := balance.New(config.Cfg.Middleware, proxy)
	if ok {
		proxy = balanceProxy
	}

	// 中间件组合
	handler := mix.New(proxy, NewGinEngine())

	if config.Cfg.Middleware != nil {
		// recovery
		handler = recovery.New(handler)
		// logger
		handler = logger.New(handler)
		// rate limit
		if config.Cfg.Middleware.RateLimit != nil {
			handler = ratelimit.New(handler)
		}
	}

	return cros.New(handler, cros.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return regexp.MustCompile(origin).MatchString(r.RequestURI)
		},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete,
			http.MethodPut, http.MethodPatch, http.MethodHead},
		AllowedHeaders: []string{"*"},
		MaxAge:         24 * 60 * 60,
	}), nil
}
