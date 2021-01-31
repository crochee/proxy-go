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

	"proxy-go/logger"
	"proxy-go/middlewares/dynamic"
	"proxy-go/util"
)

const (
	xForwardedProto             = "X-Forwarded-Proto"
	xForwardedFor               = "X-Forwarded-For"
	xForwardedHost              = "X-Forwarded-Host"
	xForwardedPort              = "X-Forwarded-Port"
	xForwardedServer            = "X-Forwarded-Server"
	xForwardedURI               = "X-Forwarded-Uri"
	xForwardedMethod            = "X-Forwarded-Method"
	xForwardedTLSClientCert     = "X-Forwarded-Tls-Client-Cert"
	xForwardedTLSClientCertInfo = "X-Forwarded-Tls-Client-Cert-Info"
	xRealIP                     = "X-Real-Ip"
	connection                  = "Connection"
	upgrade                     = "Upgrade"
)

// replaceHost is a middleware used to replace host to an URL request.
type replaceHost struct {
	next     http.Handler
	cache    map[string]*dynamic.Host
	ctx      context.Context
	hostName string
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, hostList []*dynamic.ReplaceHost) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	rh := &replaceHost{
		cache:    make(map[string]*dynamic.Host),
		next:     next,
		ctx:      ctx,
		hostName: hostname,
	}
	for _, host := range hostList {
		if host.Scheme == "" {
			host.Scheme = "http"
		}
		rh.cache[host.Name] = host.Host
	}

	return rh
}

func (r *replaceHost) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) < 2 {
		text := http.StatusText(http.StatusServiceUnavailable)
		writer.WriteHeader(http.StatusServiceUnavailable)
		if _, err := writer.Write(util.Bytes(text)); err != nil {
			logger.FromContext(r.ctx).Errorf("Error %v while writing status code", err)
		}
		return
	}
	serverName := list[1]
	if server, ok := r.cache[serverName]; ok {
		request.URL.Path = util.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, "/"+serverName))
		if request.URL.RawPath != "" {
			request.URL.RawPath = util.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, "/"+serverName))
		}

		r.rewrite(request)

		request.Header.Add(xForwardedHost, request.Host)

		request.URL.Scheme = server.Scheme
		request.URL.Host = server.Host
	}

	r.next.ServeHTTP(writer, request)
}

func (r *replaceHost) rewrite(request *http.Request) {
	if clientIP, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		clientIP = removeIPv6Zone(clientIP)

		if request.Header.Get(xRealIP) == "" {
			request.Header.Set(xRealIP, clientIP)
		}
	}

	if request.Header.Get(xForwardedProto) == "" {
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
		request.Header.Set(xForwardedProto, proto)
	}

	if xfPort := request.Header.Get(xForwardedPort); xfPort == "" {
		request.Header.Set(xForwardedPort, forwardedPort(request))
	}

	if xfHost := request.Header.Get(xForwardedHost); xfHost == "" && request.Host != "" {
		request.Header.Set(xForwardedHost, request.Host)
	}

	request.Header.Set(xForwardedServer, r.hostName)
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
	return containsHeader(connection, "upgrade") && containsHeader(upgrade, "websocket")
}

func forwardedPort(req *http.Request) string {
	if req == nil {
		return ""
	}

	if _, port, err := net.SplitHostPort(req.Host); err == nil && port != "" {
		return port
	}

	if req.Header.Get(xForwardedProto) == "https" || req.Header.Get(xForwardedProto) == "wss" {
		return "443"
	}

	if req.TLS != nil {
		return "443"
	}

	return "80"
}
