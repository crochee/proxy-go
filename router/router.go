// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"
	"time"

	"proxy-go/middlewares/dynamic"
	"proxy-go/middlewares/logger"
	"proxy-go/middlewares/ratelimit"
	"proxy-go/middlewares/recovery"
	"proxy-go/service"
)

func Route(ctx context.Context) (http.Handler, error) {
	proxy := service.NewProxyBuilder(ctx)

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
	handler = ratelimit.New(ctx, handler, &dynamic.RateLimit{
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
	if strings.HasPrefix(request.URL.Path, ProxyPrefix) {
		request.URL.Path = ensureLeadingSlash(strings.TrimPrefix(request.URL.Path, ProxyPrefix))
		if request.URL.RawPath != "" {
			request.URL.RawPath = ensureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, ProxyPrefix))
		}
		m.Proxy.ServeHTTP(writer, request)
		return
	}
	m.Proxy.ServeHTTP(writer, request)
}

func ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
