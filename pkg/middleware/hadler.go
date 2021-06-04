// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/6/5

package middleware

import (
	"context"
	"net/http"
	"sort"
	"sync/atomic"
)

var GHandlerChan = make(chan handlerList)

type Handler interface {
	Name() string
	Level() int
	Next(http.Handler)
	http.Handler
}

// Register
func Register(handlers ...Handler) {
	if len(handlers) > 0 {
		GHandlerChan <- handlers
	}
}

// Load
func Load(ctx context.Context) http.Handler {
	var httpHandler handler
	go func() {
		for {
			select {
			case handlers, ok := <-GHandlerChan:
				if ok && len(handlers) > 0 {
					sort.Sort(handlers)
					if handlers[len(handlers)-1].Level() != 0 {
						continue
					}
					temp := handlers[0]
					for index := range handlers {
						if index == 0 {
							continue
						}
						temp.Next(handlers[index])
					}
					httpHandler.value.Store(temp)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return &httpHandler
}

type handler struct {
	value atomic.Value
}

func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	value := h.value.Load()
	srv, ok := value.(http.Handler)
	if !ok {
		http.NotFound(rw, req)
		return
	}
	srv.ServeHTTP(rw, req)
}

type handlerList []Handler

func (h handlerList) Len() int {
	return len(h)
}

func (h handlerList) Less(i, j int) bool {
	return h[i].Level() > h[j].Level()
}

func (h handlerList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
