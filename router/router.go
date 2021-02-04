// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	"proxy-go/internal"
	"proxy-go/middlewares/logger"
	"proxy-go/middlewares/ratelimit"
	"proxy-go/middlewares/recovery"
	"proxy-go/middlewares/replacehost"
	"proxy-go/model"
	"proxy-go/service"
)

func Route(ctx context.Context) (http.Handler, error) {
	proxy := service.NewProxyBuilder(ctx)

	proxy = replacehost.New(ctx, proxy, []*model.ReplaceHost{
		{
			Name: "obs",
			Host: &model.Host{
				Scheme: "http",
				Host:   "localhost:8150",
			},
		},
		{
			Name: "console",
			Host: &model.Host{
				Scheme: "http",
				Host:   "localhost:8088",
			},
		},
	})
	// 中间件组合
	var (
		handler http.Handler
	)
	handler = &MixHandler{
		Proxy: proxy,
		Gin:   NewGinEngine(),
	}

	// logger
	handler = logger.New(ctx, handler)

	// recovery
	handler = recovery.New(ctx, handler)

	// rate limit
	handler = ratelimit.New(ctx, handler, &model.RateLimit{
		Every: 10 * time.Microsecond,
		Burst: 1,
	})
	return handler, nil
}

const ProxyPrefix = "proxy"

type MixHandler struct {
	Proxy http.Handler
	Gin   http.Handler
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
			m.Proxy.ServeHTTP(writer, request)
			return
		}
	}
	m.Gin.ServeHTTP(writer, request)
}
