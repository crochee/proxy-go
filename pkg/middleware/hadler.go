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
	Next(Handler) Handler
	http.Handler
}

// Register
func Register(ctx context.Context, handlers ...Handler) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case GHandlerChan <- handlers:
		return nil
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
					if handlers[0].Level() != 0 {
						continue
					}
					temp := handlers[0]
					for index := 1; index < len(handlers); index++ {
						temp = handlers[index].Next(temp)
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
	value, ok := h.value.Load().(Handler)
	if !ok {
		http.NotFound(rw, req)
		return
	}
	value.ServeHTTP(rw, req)
}

type handlerList []Handler

func (h handlerList) Len() int {
	return len(h)
}

func (h handlerList) Less(i, j int) bool {
	return h[i].Level() < h[j].Level()
}

func (h handlerList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
