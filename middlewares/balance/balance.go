// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/5

package balance

import (
	"container/heap"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"

	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/model"
)

type Balancer struct {
	ctx         context.Context
	mutex       sync.Mutex
	handlers    []*model.NamedHandler
	curDeadline float64
	hostName    string
}

func New(ctx context.Context) *Balancer {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	return &Balancer{
		ctx:      ctx,
		handlers: make([]*model.NamedHandler, 0, 4),
		hostName: hostname,
	}
}

func (b *Balancer) Len() int {
	return len(b.handlers)
}

func (b *Balancer) Less(i, j int) bool {
	return b.handlers[i].Deadline < b.handlers[j].Deadline
}

func (b *Balancer) Swap(i, j int) {
	b.handlers[i], b.handlers[j] = b.handlers[j], b.handlers[i]
}

func (b *Balancer) Push(x interface{}) {
	h, ok := x.(*model.NamedHandler)
	if !ok {
		return
	}
	b.handlers = append(b.handlers, h)
}

func (b *Balancer) Pop() interface{} {
	if b.Len() < 1 {
		return nil
	}
	h := b.handlers[len(b.handlers)-1]
	b.handlers = b.handlers[0 : len(b.handlers)-1]
	return h
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
	request.URL.Host = server.Host.Host
	server.ServeHTTP(writer, request)
}

func (b *Balancer) Update(add bool, handler *model.NamedHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for index, srvHandler := range b.handlers {
		if reflect.DeepEqual(srvHandler.Host, handler.Host) {
			if !add {
				if index == b.Len()-1 {
					b.handlers = b.handlers[:index]
					return
				}
				b.handlers = append(b.handlers[:index], b.handlers[index+1:]...)
				return
			}
			if handler.Weight > 0 {
				srvHandler.Weight = handler.Weight
			}
			return
		}
	}
	if handler.Weight <= 0 {
		return
	}
	handler.Deadline = b.curDeadline + 1/handler.Weight
	heap.Push(b, handler)
}

func (b *Balancer) nextServer() (*model.NamedHandler, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.handlers) == 0 {
		return nil, fmt.Errorf("no servers in the pool")
	}

	// Pick handler with closest deadline.
	handler, ok := heap.Pop(b).(*model.NamedHandler)
	if !ok {
		return nil, fmt.Errorf("no servers in the pool")
	}
	// todo 这个负载均衡策略待改进
	// curDeadline should be handler's deadline so that new added entry would have a fair competition environment with the old ones.
	b.curDeadline = handler.Deadline
	handler.Deadline += 1 / handler.Weight

	heap.Push(b, handler)

	logger.FromContext(b.ctx).Debugf("Service selected by WRR: %+v", handler)
	return handler, nil
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
