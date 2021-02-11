// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/5

package balance

import (
	"net"
	"net/http"
	"os"
	"strings"

	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/middlewares"
)

type Balancer struct {
	next     http.Handler
	selector Selector
	hostName string
}

func New(selector Selector, next http.Handler) *Balancer {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	return &Balancer{
		next:     next,
		selector: selector,
		hostName: hostname,
	}
}

func (b *Balancer) Name() middlewares.HandlerName {
	return middlewares.LoadBalancer
}

func (b *Balancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	node, err := b.nextNode()
	if err != nil {
		logger.FromContext(request.Context()).Warnf("get next node failed.Error:%v", err)
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	b.rewrite(request)

	request.Header.Add(internal.XForwardedHost, request.Host)
	request.URL.Scheme = node.Scheme
	request.URL.Host = node.Host

	b.next.ServeHTTP(writer, request)
}

func (b *Balancer) Update(add bool, node *Node) {
	if add && node.Weight <= 0 {
		return
	}
	b.selector.Update(add, node)
}

func (b *Balancer) nextNode() (*Node, error) {
	return b.selector.Next()
}

func (b *Balancer) NodeList() []*Node {
	return b.selector.List()
}

func (b *Balancer) rewrite(request *http.Request) {
	if clientIP, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		clientIP = removeIPv6Zone(clientIP)

		if request.Header.Get(internal.XRealIP) == "" {
			request.Header.Set(internal.XRealIP, clientIP)
		}
	}

	if request.Header.Get(internal.XForwardedProto) == "" {
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
		request.Header.Set(internal.XForwardedProto, proto)
	}

	if xfPort := request.Header.Get(internal.XForwardedPort); xfPort == "" {
		request.Header.Set(internal.XForwardedPort, forwardedPort(request))
	}

	if xfHost := request.Header.Get(internal.XForwardedHost); xfHost == "" && request.Host != "" {
		request.Header.Set(internal.XForwardedHost, request.Host)
	}

	request.Header.Set(internal.XForwardedServer, b.hostName)
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
	return containsHeader(internal.Connection, "upgrade") && containsHeader(internal.Upgrade, "websocket")
}

func forwardedPort(req *http.Request) string {
	if req == nil {
		return ""
	}

	if _, port, err := net.SplitHostPort(req.Host); err == nil && port != "" {
		return port
	}

	if req.Header.Get(internal.XForwardedProto) == "https" || req.Header.Get(internal.XForwardedProto) == "wss" {
		return "443"
	}

	if req.TLS != nil {
		return "443"
	}

	return "80"
}
