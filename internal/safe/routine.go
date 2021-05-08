// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package safe

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/crochee/proxy-go/logger"
)

type pool struct {
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a Pool.
func NewPool(parentCtx context.Context) *pool {
	ctx, cancel := context.WithCancel(parentCtx)
	return &pool{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Go starts a recoverable goroutine with a context.
func (p *pool) Go(goroutine func(ctx context.Context)) {
	p.waitGroup.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.FromContext(p.ctx).Errorf("[Recovery] panic happened.Error:%v\n.Stack:\n%s", debug.Stack())
			}
			p.waitGroup.Done()
		}()
		goroutine(p.ctx)
	}()
}

// Stop stops all started routines, waiting for their termination.
func (p *pool) Stop() {
	p.cancel()
	p.waitGroup.Wait()
}
