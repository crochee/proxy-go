// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package logger

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"proxy-go/logger"
)

type loggerHandler struct {
	ctx  context.Context
	next http.Handler
}

func New(ctx context.Context, next http.Handler) (http.Handler, error) {
	return &loggerHandler{
		ctx:  ctx,
		next: next,
	}, nil
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

	param.Now = time.Now().Local()
	param.Last = param.Now.Sub(start)
	if param.Last > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Last = param.Last - param.Last%time.Second
	}
	param.Status = crw.Status()
	param.Size = crw.Size()

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
