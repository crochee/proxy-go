// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/27

// Package prometheusx
package metrics

import (
	"github.com/crochee/proxy-go/version"
	"github.com/prometheus/client_golang/prometheus"
	"sync/atomic"
)

var (
	Enable              atomic.Value
	ReqDurHistogramVec  *prometheus.HistogramVec
	ReqCodeTotalCounter *prometheus.CounterVec
)

// DefineMetrics init metrics
func DefineMetrics() {
	Enable.Store(true)
	ReqDurHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: version.ServiceName,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	}, []string{"protocol", "method", "path", "code"})
	ReqCodeTotalCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: version.ServiceName,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "http server requests error count.",
	}, []string{"protocol", "method", "path", "code"})
}
