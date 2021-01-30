// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"

	"proxy-go/config"
	"proxy-go/middlewares"
)

func Route(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	proxy, err := middlewares.BuildProxy(0, *cfg.Server.Medata[0])
	if err != nil {
		return nil, err
	}
	var proxyHandler http.Handler
	if proxy, err = middlewares.NewRecovery(ctx, proxy); err != nil {
		return nil, err
	}

	return &MixHandler{
		Proxy: proxyHandler,
		Gin:   GinRun(),
	}, nil
}

const ProxyPrefix = "proxy"

type MixHandler struct {
	Proxy http.Handler
	Gin   http.Handler
}

func (m *MixHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if strings.HasPrefix(request.URL.Path, ProxyPrefix) {
		m.Proxy.ServeHTTP(writer, request)
		return
	}
	m.Proxy.ServeHTTP(writer, request)
}
