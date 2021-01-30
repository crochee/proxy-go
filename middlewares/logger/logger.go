// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package logger

import (
	"context"
	"net"
	"net/http"
	"proxy-go/logger"
	"strings"
	"time"
)

type loggerHandler struct {
	ctx  context.Context
	next http.Handler
}

func NewLogger(ctx context.Context, next http.Handler) (http.Handler, error) {
	return &loggerHandler{
		ctx:  ctx,
		next: next,
	}, nil
}

func (l *loggerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now().Local()
	path := request.URL.Path
	raw := request.URL.RawQuery
	scheme := "http"
	proto := request.Proto
	if request.TLS != nil {
		scheme = "https"
	}
	if raw != "" {
		path = path + "?" + raw
	}
	crw := newCaptureResponseWriter(writer)

	l.next.ServeHTTP(crw, request)

	latency := time.Now().Local().Sub(start)
	if latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		latency = latency - latency%time.Second
	}

	logger.FromContext(l.ctx).Infof(
		"[PROXY] %v | %3d | %13v | %15s | %-7s | %5s | %10s |%8d| %#v",
		time.Now().Local().Format("2006/01/02 - 15:04:05"),
		crw.Status(),
		latency,
		scheme,
		proto,
		clientIp(request),
		request.Method,
		crw.Size(),
		path,
	)
}

func clientIp(request *http.Request) string {
	clientIP := request.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(request.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr)); err == nil {
		return ip
	}

	return "-"
}
