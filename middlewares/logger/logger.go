// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package logger

import (
	"context"
	"net"
	"net/http"
	"proxy-go/internal"
	"strings"
	"time"

	"proxy-go/logger"
	"proxy-go/middlewares"
)

type loggerHandler struct {
	ctx  context.Context
	next http.Handler
}

func New(ctx context.Context, next http.Handler) middlewares.MiddleWare {
	return &loggerHandler{
		ctx:  ctx,
		next: next,
	}
}

func (l *loggerHandler) Name() string {
	return "Logger"
}

func (l *loggerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now().Local()
	param := &LogFormatterParams{
		Scheme:   "HTTP",
		Proto:    request.Proto,
		ClientIp: clientIp(request),
		Method:   request.Method,
		Path:     request.URL.Path,
	}
	if request.TLS != nil {
		param.Scheme = "HTTPS"
	}
	if request.URL.RawQuery != "" {
		param.Path += "?" + request.URL.RawQuery
	}

	crw := newCaptureResponseWriter(writer)

	l.next.ServeHTTP(crw, request)

	param.Status = crw.Status()
	param.Size = crw.Size()
	param.Now = time.Now().Local()
	param.Last = param.Now.Sub(start)
	if param.Last > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Last = param.Last - param.Last%time.Second
	}
	logger.FromContext(l.ctx).Infof(
		"[PROXY] %v | %3d | %13v | %15s | %-7s | %5s | %10s |%8d| %#v",
		param.Now.Format("2006/01/02 - 15:04:05"),
		param.Status,
		param.Last,
		param.Scheme,
		param.Proto,
		param.ClientIp,
		param.Method,
		param.Size,
		param.Path,
	)
}

type LogFormatterParams struct {
	Now      time.Time
	Status   int
	Last     time.Duration
	Scheme   string
	Proto    string
	ClientIp string
	Method   string
	Size     int64
	Path     string
}

func clientIp(request *http.Request) string {
	clientIP := request.Header.Get(internal.XForwardedFor)
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(request.Header.Get(internal.XRealIP))
	}
	if clientIP != "" {
		return clientIP
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr)); err == nil {
		return ip
	}

	return "-"
}
