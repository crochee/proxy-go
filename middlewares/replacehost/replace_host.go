// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package replacehost

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"proxy-go/internal"
	"proxy-go/model"
)

// replaceHost is a middleware used to replace host to an URL request.
type replaceHost struct {
	next     http.Handler
	cache    *sync.Map
	ctx      context.Context
	hostName string
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, hostList []*model.ReplaceHost) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	rh := &replaceHost{
		cache:    new(sync.Map),
		next:     next,
		ctx:      ctx,
		hostName: hostname,
	}
	for _, host := range hostList {
		if host.Scheme == "" {
			host.Scheme = "http"
		}
		rh.cache.Store(host.Name, host.Host)
	}

	return rh
}

func (r *replaceHost) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) > 1 {
		serverName := list[1]
		if server, ok := r.cache.Load(serverName); ok {
			if value, ok := server.(*model.Host); ok {

				request.URL.Path = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, "/"+serverName))
				if request.URL.RawPath != "" {
					request.URL.RawPath = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, "/"+serverName))
				}

				r.rewrite(request)

				request.Header.Add(model.XForwardedHost, request.Host)

				request.URL.Scheme = value.Scheme
				request.URL.Host = value.Host
			}
		}
	}
	r.next.ServeHTTP(writer, request)
}

func (r *replaceHost) rewrite(request *http.Request) {
	if clientIP, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		clientIP = removeIPv6Zone(clientIP)

		if request.Header.Get(model.XRealIP) == "" {
			request.Header.Set(model.XRealIP, clientIP)
		}
	}

	if request.Header.Get(model.XForwardedProto) == "" {
		var proto string
		if isWebsocketRequest(request) {
			if request.TLS != nil {
				proto = "wss"
			} else {
				proto = "ws"
			}
		} else {
			if request.TLS != nil {
				proto = "https"
			} else {
				proto = "http"
			}
		}
		request.Header.Set(model.XForwardedProto, proto)
	}

	if xfPort := request.Header.Get(model.XForwardedPort); xfPort == "" {
		request.Header.Set(model.XForwardedPort, forwardedPort(request))
	}

	if xfHost := request.Header.Get(model.XForwardedHost); xfHost == "" && request.Host != "" {
		request.Header.Set(model.XForwardedHost, request.Host)
	}

	request.Header.Set(model.XForwardedServer, r.hostName)
}

// removeIPv6Zone removes the zone if the given IP is an ipv6 address and it has {zone} information in it,
// like "[fe80::d806:a55d:eb1b:49cc%vEthernet (vmxnet3 Ethernet Adapter - Virtual Switch)]:64692".
func removeIPv6Zone(clientIP string) string {
	return strings.Split(clientIP, "%")[0]
}

// isWebsocketRequest returns whether the specified HTTP request is a websocket handshake request.
func isWebsocketRequest(req *http.Request) bool {
	containsHeader := func(name, value string) bool {
		items := strings.Split(req.Header.Get(name), ",")
		for _, item := range items {
			if value == strings.ToLower(strings.TrimSpace(item)) {
				return true
			}
		}
		return false
	}
	return containsHeader(model.Connection, "upgrade") && containsHeader(model.Upgrade, "websocket")
}

func forwardedPort(req *http.Request) string {
	if req == nil {
		return ""
	}

	if _, port, err := net.SplitHostPort(req.Host); err == nil && port != "" {
		return port
	}

	if req.Header.Get(model.XForwardedProto) == "https" || req.Header.Get(model.XForwardedProto) == "wss" {
		return "443"
	}

	if req.TLS != nil {
		return "443"
	}

	return "80"
}
