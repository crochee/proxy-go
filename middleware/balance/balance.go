// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/5

package balance

import (
	"net"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/internal/selector"
	"github.com/crochee/proxy-go/logger"
)

type SelectorInfo struct {
	*dynamic.Balance
	selector.Selector
}

type Balancer struct {
	next http.Handler
	// path method service
	serviceApi map[string]map[string]string
	// service selector
	NameSelector map[string]*SelectorInfo
	rw           sync.RWMutex
	hostName     string
}

func New(cfg *dynamic.BalanceConfig, next http.Handler) http.Handler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	b := &Balancer{
		next:         next,
		serviceApi:   make(map[string]map[string]string),
		NameSelector: make(map[string]*SelectorInfo),
		hostName:     hostname,
	}
	for _, api := range cfg.RegisterApis {
		if _, ok := b.serviceApi[api.Path]; !ok {
			b.serviceApi[api.Path] = make(map[string]string)
		}
		b.serviceApi[api.Path][api.Method] = api.ServiceName
	}
	for _, balance := range cfg.Transfers {
		b.NameSelector[balance.ServiceName] = &SelectorInfo{
			Balance:  &balance.Balance,
			Selector: createSelector(&balance.Balance),
		}
	}
	return b
}

func (b *Balancer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	apis, ok := b.serviceApi[request.URL.Path]
	if !ok {
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	var serviceName string
	if serviceName, ok = apis[request.Method]; !ok {
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	b.rw.RLock()
	s, ok := b.NameSelector[serviceName]
	b.rw.RUnlock()
	if !ok {
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	node, err := s.Next()
	if err != nil {
		logger.FromContext(request.Context()).Errorf("get next node failed.Error:%v", err)
		http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}

	b.rewrite(request)

	request.Header.Add(internal.XForwardedHost, request.Host)
	request.URL.Scheme = node.Scheme
	request.URL.Host = node.Host

	b.next.ServeHTTP(writer, request)
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

func createSelector(balance *dynamic.Balance) selector.Selector {
	var s selector.Selector
	switch strings.Title(balance.Selector) {
	case "Random":
		s = selector.NewRandom()
	case "RoundRobin":
		s = selector.NewRoundRobin()
	case "Heap":
		s = selector.NewHeap()
	case "Wrr":
		fallthrough
	default:
		s = selector.NewWeightRoundRobin()
	}
	for _, node := range balance.Nodes {
		s.AddNode(&selector.Node{
			Scheme:   node.Scheme,
			Host:     node.Host,
			Metadata: node.Metadata,
			Weight:   node.Weight,
		})
	}
	return s
}
