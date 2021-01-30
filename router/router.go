// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"context"
	"net/http"
	"strings"

	"proxy-go/middlewares"
	"proxy-go/service"
)

func Route(ctx context.Context) (http.Handler, error) {
	proxy, err := service.NewProxyBuilder(0)
	if err != nil {
		return nil, err
	}
	var proxyHandler http.Handler
	if proxy, err = middlewares.NewRecovery(ctx, proxy); err != nil {
		return nil, err
	}

	return &MixHandler{
		Proxy: proxyHandler,
		Gin:   NewGinEngine(),
	}, nil
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
