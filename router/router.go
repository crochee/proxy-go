// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	"proxy-go/config/dynamic"
	"proxy-go/internal"
	"proxy-go/middlewares/balance"
	"proxy-go/middlewares/logger"
	"proxy-go/middlewares/ratelimit"
	"proxy-go/middlewares/recovery"
	"proxy-go/middlewares/selecthandler"
	"proxy-go/server"
	"proxy-go/service"
)

func ChainBuilder(ctx context.Context, watcher *server.Watcher) (http.Handler, error) {
	proxy := service.NewProxyBuilder(ctx)

	balancer := balance.New(ctx, balance.NewRandom(), proxy)

	watcher.AddListener(balancer.Name(), func(config *dynamic.Config) {
		balancer.Update(true, &balance.Node{
			Scheme: "http",
			Host:   "127.0.0.1:8150",
		}, 1)
	})

	switchHandler := selecthandler.New(ctx)

	watcher.AddListener(switchHandler.Name(), func(config *dynamic.Config) {
		switchHandler.Update("obs", balancer)
	})

	// 中间件组合
	var (
		handler http.Handler
	)
	handler = &MixHandler{
		proxy: switchHandler,
		api:   NewGinEngine(),
	}

	// recovery
	handler = recovery.New(ctx, handler)

	// logger
	handler = logger.New(ctx, handler)

	// rate limit
	handler = ratelimit.New(ctx, handler, &dynamic.RateLimit{
		Every: 10 * time.Microsecond,
		Burst: 1,
	})
	return handler, nil
}

const ProxyPrefix = "proxy"

type MixHandler struct {
	proxy http.Handler
	api   http.Handler
}

func (m *MixHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) > 1 {
		if list[1] == ProxyPrefix {
			prefix := "/" + ProxyPrefix
			request.URL.Path = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, prefix))
			if request.URL.RawPath != "" {
				request.URL.RawPath = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, prefix))
			}
			m.proxy.ServeHTTP(writer, request)
			return
		}
	}
	m.api.ServeHTTP(writer, request)
}
