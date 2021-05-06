// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package logger

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/middleware"
)

type loggerHandler struct {
	next http.Handler
}

func New(next http.Handler) middleware.Handler {
	return &loggerHandler{
		next: next,
	}
}

func (l *loggerHandler) NameSpace() string {
	return "AccessLog"
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
	var buf strings.Builder
	if request.URL.RawQuery != "" {
		buf.WriteString(param.Path)
		buf.WriteByte('?')
		buf.WriteString(request.URL.RawQuery)
		param.Path = buf.String()
		buf.Reset()
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
	buf.WriteString("[PROXY] ")
	buf.WriteString(param.Now.Format("2006/01/02 - 15:04:05"))
	buf.WriteString(" | ")
	buf.WriteString(strconv.Itoa(param.Status))
	buf.WriteString(" | ")
	buf.WriteString(param.Last.String())
	buf.WriteString(" | ")
	buf.WriteString(param.Scheme)
	buf.WriteString(" | ")
	buf.WriteString(param.Proto)
	buf.WriteString(" | ")
	buf.WriteString(param.ClientIp)
	buf.WriteString(" |")
	buf.WriteString(param.Method)
	buf.WriteString("| ")
	buf.WriteString(strconv.Itoa(int(param.Size)))
	buf.WriteString(" | ")
	buf.WriteString(param.Path)
	logger.FromContext(request.Context()).Info(buf.String())
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
