// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/5

package balance

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/model"
)

type Balancer struct {
	ctx      context.Context
	selector Selector
	hostName string
}

func New(ctx context.Context, selector Selector) *Balancer {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	return &Balancer{
		ctx:      ctx,
		selector: selector,
		hostName: hostname,
	}
}

func (b *Balancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	server, err := b.nextServer()
	if err != nil {
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable)+err.Error(),
			http.StatusServiceUnavailable)
		return
	}

	b.rewrite(request)

	request.Header.Add(model.XForwardedHost, request.Host)
	request.URL.Scheme = server.Scheme
	request.URL.Host = server.Host
	server.ServeHTTP(writer, request)
}

func (b *Balancer) Update(add bool, handler *model.NamedHandler) {
	if add && handler.Weight <= 0 {
		logger.FromContext(b.ctx).Warnf("add handler failed.it's Weight is %f", handler.Weight)
		return
	}
	b.selector.Update(add, handler, handler.Weight)
}

func (b *Balancer) nextServer() (*model.NamedHandler, error) {
	handler, err := b.selector.Next()
	if err != nil {
		return nil, err
	}
	srv, ok := handler.(*model.NamedHandler)
	if !ok {
		return nil, fmt.Errorf("no servers in the pool")
	}
	return srv, nil
}

func (b *Balancer) rewrite(request *http.Request) {
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

	request.Header.Set(model.XForwardedServer, b.hostName)
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
