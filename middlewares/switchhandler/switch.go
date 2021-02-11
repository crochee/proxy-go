// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/31

package switchhandler

import (
	"net/http"
	"strings"
	"sync"

	"proxy-go/internal"
	"proxy-go/middlewares"
)

// SwitchHandler is a middleware used to switch handler.
type SwitchHandler struct {
	cache *sync.Map
}

// New creates a new handler.
func New() *SwitchHandler {
	return &SwitchHandler{
		cache: new(sync.Map),
	}
}

func (s *SwitchHandler) Name() middlewares.HandlerName {
	return middlewares.Switcher
}

func (s *SwitchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) > 1 {
		serverName := list[1]
		if server, ok := s.cache.Load(serverName); ok {
			if value, ok := server.(http.Handler); ok {

				request.URL.Path = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, "/"+serverName))
				if request.URL.RawPath != "" {
					request.URL.RawPath = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, "/"+serverName))
				}

				value.ServeHTTP(writer, request)
				return
			}
		}
	}
	http.Error(writer, internal.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

func (s *SwitchHandler) Store(name string, handler http.Handler) {
	s.cache.Store(name, handler)
}

func (s *SwitchHandler) Load(serviceName string) (http.Handler, bool) {
	value, ok := s.cache.Load(serviceName)
	if !ok {
		return nil, false
	}
	return value.(http.Handler), true
}

func (s *SwitchHandler) Delete(serviceName string) {
	s.cache.Delete(serviceName)
}

func (s *SwitchHandler) Range(function func(key, value interface{}) bool) {
	s.cache.Range(function)
}
