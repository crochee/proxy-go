// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package router

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal"
	"github.com/crochee/proxy-go/middleware"
	"github.com/crochee/proxy-go/middleware/balance"
	"github.com/crochee/proxy-go/middleware/cros"
	"github.com/crochee/proxy-go/middleware/logger"
	"github.com/crochee/proxy-go/middleware/ratelimit"
	"github.com/crochee/proxy-go/middleware/recovery"
	"github.com/crochee/proxy-go/middleware/switchhandler"
	"github.com/crochee/proxy-go/service"
	"github.com/crochee/proxy-go/service/server"
)

func ChainBuilder(watcher *server.Watcher) (http.Handler, error) {
	proxy := service.NewProxyBuilder()

	switchHandler := switchhandler.New()

	watcher.AddListener(
		middleware.CompleteAction(switchHandler.Name(), middleware.Update),
		func(config *dynamic.Config, _ chan<- interface{}) {
			if !config.Switcher.Add {
				switchHandler.Delete(config.Switcher.ServiceName)
				return
			}
			var balancer *balance.Balancer
			handler, ok := switchHandler.Load(config.Switcher.ServiceName)
			if !ok {
				balancer = balance.New(balance.NewRoundRobin(), proxy)
				switchHandler.Store(config.Switcher.ServiceName, balancer)
			} else {
				if balancer, ok = handler.(*balance.Balancer); !ok {
					switchHandler.Delete(config.Switcher.ServiceName)
					return
				}
			}
			balancer.Update(config.Switcher.Node.Add, &balance.Node{
				Scheme:   config.Switcher.Node.Scheme,
				Host:     config.Switcher.Node.Host,
				Metadata: config.Switcher.Node.Metadata,
				Weight:   config.Switcher.Node.Weight,
			})
		})

	watcher.AddListener(
		middleware.CompleteAction(switchHandler.Name(), middleware.List),
		func(config *dynamic.Config, response chan<- interface{}) {
			list := make([]*dynamic.Switch, 0, 4)
			switchHandler.Range(func(key, value interface{}) bool {
				keyStr, ok := key.(string)
				if !ok {
					return true
				}
				var node *balance.Balancer
				if node, ok = value.(*balance.Balancer); !ok {
					return true
				}
				nodeList := node.NodeList()
				if len(nodeList) == 0 {
					list = append(list, &dynamic.Switch{
						ServiceName: keyStr,
					})
					return true
				}
				for _, nodeValue := range nodeList {
					list = append(list, &dynamic.Switch{
						ServiceName: keyStr,
						Node: dynamic.BalanceNode{
							Scheme:   nodeValue.Scheme,
							Host:     nodeValue.Host,
							Metadata: nodeValue.Metadata,
							Weight:   nodeValue.Weight,
						},
					})
				}
				return true
			})
			response <- list
		})

	// 中间件组合
	var (
		handler http.Handler
	)
	handler = &MixHandler{
		proxy: switchHandler,
		api:   NewGinEngine(),
	}

	// recovery
	handler = recovery.New(handler)

	// logger
	handler = logger.New(handler)

	limit := ratelimit.New(handler)

	watcher.AddListener(
		middleware.CompleteAction(limit.Name(), middleware.Update),
		func(config *dynamic.Config, _ chan<- interface{}) {
			limit.Update(config.Limit)
		})

	watcher.AddListener(
		middleware.CompleteAction(limit.Name(), middleware.Get),
		func(config *dynamic.Config, response chan<- interface{}) {
			response <- limit.Get()
		})
	return cros.New(limit, cros.Options{
		AllowOriginRequestFunc: func(r *http.Request, origin string) bool {
			return regexp.MustCompile(origin).MatchString(r.RequestURI)
		},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete,
			http.MethodPut, http.MethodPatch, http.MethodHead},
		AllowedHeaders: []string{"*"},
		MaxAge:         24 * 60 * 60,
	}), nil
}

const ProxyPrefix = "proxy"

type MixHandler struct {
	proxy http.Handler
	api   http.Handler
}

func (m *MixHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	list := strings.SplitN(request.URL.Path, "/", 3)
	if len(list) > 1 {
		if list[1] == ProxyPrefix {
			prefix := "/" + ProxyPrefix
			request.URL.Path = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.Path, prefix))
			if request.URL.RawPath != "" {
				request.URL.RawPath = internal.EnsureLeadingSlash(strings.TrimPrefix(request.URL.RawPath, prefix))
			}
			m.proxy.ServeHTTP(writer, request)
			return
		}
	}
	m.api.ServeHTTP(writer, request)
}
