// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"

	"proxy-go/config/dynamic"
	"proxy-go/internal"
	"proxy-go/middlewares/balance"
	"proxy-go/middlewares/logger"
	"proxy-go/middlewares/recovery"
	"proxy-go/middlewares/switchhandler"
	"proxy-go/server"
	"proxy-go/service"
)

func ChainBuilder(ctx context.Context, watcher *server.Watcher) (http.Handler, error) {
	proxy := service.NewProxyBuilder(ctx)

	switchHandler := switchhandler.New(ctx)

	watcher.AddListener(switchHandler.Name(), func(config *dynamic.Config) {
		if !config.Switcher.Add {
			switchHandler.Delete(config.Switcher.ServiceName)
			return
		}
		var balancer *balance.Balancer
		handler, ok := switchHandler.Load(config.Switcher.ServiceName)
		if !ok {
			balancer = balance.New(ctx, balance.NewRoundRobin(), proxy)
			switchHandler.Store(config.Switcher.ServiceName, balancer)
		} else {
			if balancer, ok = handler.(*balance.Balancer); !ok {
				switchHandler.Delete(config.Switcher.ServiceName)
				return
			}
		}
		balancer.Update(config.Switcher.Node.Add, &balance.Node{
			Scheme:   config.Switcher.Node.Scheme,
			Host:     config.Switcher.Node.Host,
			Metadata: config.Switcher.Node.Metadata,
		}, config.Switcher.Node.Weight)
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
	//limit := ratelimit.New(ctx, handler)
	//
	//watcher.AddListener(limit.Name(), func(config *dynamic.Config) {
	//	limit.Update(config.Limit)
	//})

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
