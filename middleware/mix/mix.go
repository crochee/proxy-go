// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/7

// Package mix
package mix

import (
	"net/http"
	"strings"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/middleware"
)

const ProxyPrefix = "proxy"

type MixHandler struct {
	proxy http.Handler
	api   http.Handler
}

func New(proxy, api http.Handler) middleware.Handler {
	return &MixHandler{
		proxy: proxy,
		api:   api,
	}
}

func (m *MixHandler) NameSpace() string {
	return "mix"
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
