package accesslog

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crochee/proxy/internal"
	"github.com/crochee/proxy/internal/writer"
	"github.com/crochee/proxy/pkg/logger"
	"github.com/crochee/proxy/pkg/middleware"
)

func New(log logger.Builder) *accessLog {
	if log == nil {
		log = logger.NoLogger{}
	}
	return &accessLog{
		log: log,
	}
}

type accessLog struct {
	next middleware.Handler
	log  logger.Builder
}

func (l *accessLog) Name() string {
	return "ACCESS_LOG"
}

func (l *accessLog) Level() int {
	return 4
}

func (l *accessLog) Next(handler middleware.Handler) middleware.Handler {
	l.next = handler
	return l
}

func (l *accessLog) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now().Local()
	param := &LogFormatterParams{
		Scheme:   "HTTP",
		Proto:    req.Proto,
		ClientIp: clientIp(req),
		Method:   req.Method,
		Path:     req.URL.Path,
	}
	if req.TLS != nil {
		param.Scheme = "HTTPS"
	}
	var buf strings.Builder
	if req.URL.RawQuery != "" {
		buf.WriteString(param.Path)
		buf.WriteByte('?')
		buf.WriteString(req.URL.RawQuery)
		param.Path = buf.String()
		buf.Reset()
	}

	crw := writer.NewCaptureResponseWriter(rw)

	l.next.ServeHTTP(crw, req)

	param.Status = crw.Status()
	param.Size = crw.Size()
	param.Now = time.Now().Local()
	param.Last = param.Now.Sub(start)
	if param.Last > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Last = param.Last - param.Last%time.Second
	}
	buf.WriteString("[ACCESS_LOG]")
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
	l.log.Info(buf.String())
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
