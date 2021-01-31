// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package replacehost

import (
	"context"
	"net/http"
	"strings"

	"proxy-go/logger"
	"proxy-go/middlewares/dynamic"
	"proxy-go/util"
)

const ReplacedHostHeader = "X-Replaced-Host"

// replaceHost is a middleware used to replace host to an URL request.
type replaceHost struct {
	next  http.Handler
	cache map[string]*dynamic.Host
	ctx   context.Context
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, hostList []*dynamic.ReplaceHost) http.Handler {
	rh := &replaceHost{
		cache: make(map[string]*dynamic.Host),
		next:  next,
		ctx:   ctx,
	}
	for _, host := range hostList {
		if host.Scheme == "" {
			host.Scheme = "http"
		}
		rh.cache[host.Name] = host.Host
	}

	return rh
}

func (r *replaceHost) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) < 2 {
		text := http.StatusText(http.StatusServiceUnavailable)
		writer.WriteHeader(http.StatusServiceUnavailable)
		if _, err := writer.Write(util.Bytes(text)); err != nil {
			logger.FromContext(r.ctx).Errorf("Error %v while writing status code", err)
		}
		return
	}
	serverName := list[1]
	if server, ok := r.cache[serverName]; ok {
		request.URL.Path = util.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, "/"+serverName))
		if request.URL.RawPath != "" {
			request.URL.RawPath = util.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, "/"+serverName))
		}
		request.Header.Add(ReplacedHostHeader, request.URL.Host)
		request.URL.Scheme = server.Scheme
		request.URL.Host = server.Host
		request.RequestURI = request.URL.RequestURI()
	}

	r.next.ServeHTTP(writer, request)
}
