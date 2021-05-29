// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/27

// Package metric
package metric

import (
	"github.com/crochee/proxy-go/pkg/metrics"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/crochee/proxy-go/pkg/writer"
)

type metric struct {
	next         http.Handler
	reqDur       *prometheus.HistogramVec
	reqCodeTotal *prometheus.CounterVec
}

// New create metric http.Handler
func New(next http.Handler) *metric {
	return &metric{
		next: next,
	}
}

func (m *metric) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if metrics.Enable.Load() != true {
		m.next.ServeHTTP(rw, req)
		return
	}
	var labels []string
	labels = append(labels, getRequestProtocol(req), req.Method, req.URL.Path)

	start := time.Now()

	crw := writer.NewCaptureResponseWriter(rw)

	m.next.ServeHTTP(crw, req)

	labels = append(labels, strconv.Itoa(crw.Status()))

	m.reqDur.WithLabelValues(labels...).Observe(float64(time.Since(start).Nanoseconds()))

	if crw.Status()/100 != 2 {
		m.reqCodeTotal.WithLabelValues(labels...).Inc()
	}
}

func getRequestProtocol(req *http.Request) string {
	switch {
	case isWebsocketRequest(req):
		return "websocket"
	case isSSERequest(req):
		return "sse"
	default:
		return "http"
	}
}

// isWebsocketRequest determines if the specified HTTP request is a websocket handshake request.
func isWebsocketRequest(req *http.Request) bool {
	return containsHeader(req, "Connection", "upgrade") && containsHeader(req, "Upgrade", "websocket")
}

// isSSERequest determines if the specified HTTP request is a request for an event subscription.
func isSSERequest(req *http.Request) bool {
	return containsHeader(req, "Accept", "text/event-stream")
}

func containsHeader(req *http.Request, name, value string) bool {
	items := strings.Split(req.Header.Get(name), ",")
	for _, item := range items {
		if value == strings.ToLower(strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}
