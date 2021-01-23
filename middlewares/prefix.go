// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/23

package middlewares

import (
	"context"
	"net/http"
	"strings"

	"proxy-go/config/dynamic"
)

const (
	// ForwardedPrefixHeader is the default header to set prefix.
	ForwardedPrefixHeader = "X-Forwarded-Prefix"
)

type stripPrefix struct {
	ctx        context.Context
	next       http.Handler
	prefixList []string
}

func NewRegexpPath(ctx context.Context, medata dynamic.Middleware, next http.Handler) (http.Handler, error) {
	return &stripPrefix{
		ctx:        ctx,
		next:       next,
		prefixList: medata.Prefix.Path,
	}, nil
}

func (s *stripPrefix) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	for _, prefix := range s.prefixList {
		if strings.HasPrefix(request.URL.Path, prefix) {
			request.URL.Path = s.getPrefixStripped(request.URL.Path, prefix)
			if request.URL.RawPath != "" {
				request.URL.RawPath = s.getPrefixStripped(request.URL.RawPath, prefix)
			}
			request.Header.Add(ForwardedPrefixHeader, strings.TrimSpace(prefix))
			request.RequestURI = request.URL.RequestURI()
			s.next.ServeHTTP(writer, request)
			return
		}
	}
	s.next.ServeHTTP(writer, request)
}

func (s *stripPrefix) serveRequest(rw http.ResponseWriter, req *http.Request, prefix string) {
	req.Header.Add(ForwardedPrefixHeader, prefix)
	req.RequestURI = req.URL.RequestURI()
	s.next.ServeHTTP(rw, req)
}

func (s *stripPrefix) getPrefixStripped(urlPath, prefix string) string {
	return ensureLeadingSlash(strings.TrimPrefix(urlPath, prefix))
}

func ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
