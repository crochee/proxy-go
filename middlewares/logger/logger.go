// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package logger

import (
	"net"
	"net/http"
	"strings"
	"time"

	"proxy-go/internal"
	"proxy-go/logger"
	"proxy-go/middlewares"
)

type loggerHandler struct {
	next http.Handler
}

func New(next http.Handler) middlewares.Handler {
	return &loggerHandler{
		next: next,
	}
}

func (l *loggerHandler) Name() middlewares.HandlerName {
	return middlewares.Logger
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
		buf := internal.GetBuffer()
		buf.AppendString(param.Path)
		buf.AppendByte('?')
		buf.AppendString(request.URL.RawQuery)
		param.Path = buf.String()
		buf.Free()
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
	buf := internal.GetBuffer()
	buf.AppendString("[PROXY] ")
	buf.AppendTime(param.Now, "2006/01/02 - 15:04:05")
	buf.AppendString(" | ")
	buf.AppendInt(int64(param.Status))
	buf.AppendString(" | ")
	buf.AppendString(param.Last.String())
	buf.AppendString(" | ")
	buf.AppendString(param.Scheme)
	buf.AppendString(" | ")
	buf.AppendString(param.Proto)
	buf.AppendString(" | ")
	buf.AppendString(param.ClientIp)
	buf.AppendString(" |")
	buf.AppendString(param.Method)
	buf.AppendString("| ")
	buf.AppendInt(param.Size)
	buf.AppendString(" | ")
	buf.AppendString(param.Path)
	logger.FromContext(request.Context()).Info(buf.String())
	buf.Free()
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
	if addr := request.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr)); err == nil {
		return ip
	}

	return "-"
}
