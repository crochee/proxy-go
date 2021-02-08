// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package server

import (
	"context"
	"sync"

	"proxy-go/config/dynamic"
	"proxy-go/safe"
)

type Message struct {
	Name    string
	Content *dynamic.Config
}

type Watcher struct {
	ctx      context.Context
	pool     *safe.Pool
	storeMap *sync.Map //map[string]DynamicFunc
	message  chan *Message
}

var GlobalWatcher *Watcher

func NewWatcher(ctx context.Context, pool *safe.Pool) *Watcher {
	return &Watcher{
		ctx:      ctx,
		pool:     pool,
		storeMap: new(sync.Map),
		message:  make(chan *Message, 100),
	}
}

type DynamicFunc func(*dynamic.Config)

func (w *Watcher) Start() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case message, ok := <-w.message:
			if !ok {
				continue
			}
			function, ok := w.storeMap.Load(message.Name)
			if !ok {
				continue
			}
			dynamicFunc, ok := function.(DynamicFunc)
			if !ok {
				continue
			}
			w.pool.Go(func(ctx context.Context) {
				dynamicFunc(message.Content)
			})
		}
	}
}

func (w *Watcher) AddListener(name string, function DynamicFunc) {
	w.storeMap.Store(name, function)
}

func (w *Watcher) Entry() chan<- *Message {
	return w.message
}

func (w *Watcher) Stop() {
	w.pool.Stop()
	close(w.message)
}
