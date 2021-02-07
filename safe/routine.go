// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package safe

import (
	"context"
	"runtime/debug"
	"sync"

	"proxy-go/logger"
)

type routineCtx func(ctx context.Context)

// Pool is a pool of go routines.
type Pool struct {
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewPool creates a Pool.
func NewPool(parentCtx context.Context) *Pool {
	ctx, cancel := context.WithCancel(parentCtx)
	return &Pool{
		ctx:    ctx,
		cancel: cancel,
	}
}

// GoCtx starts a recoverable goroutine with a context.
func (p *Pool) GoCtx(goroutine routineCtx) {
	p.waitGroup.Add(1)
	Go(func() {
		defer p.waitGroup.Done()
		goroutine(p.ctx)
	})
}

// Go starts a recoverable goroutine with a context.
func (p *Pool) Go(goroutine routineCtx) {
	p.waitGroup.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.FromContext(p.ctx).Errorf("Error in Go routine: %v", err)
				logger.FromContext(p.ctx).Errorf("Stack: %s", debug.Stack())
			}
			p.waitGroup.Done()
		}()
		goroutine(p.ctx)
	}()
}

// Stop stops all started routines, waiting for their termination.
func (p *Pool) Stop() {
	p.cancel()
	p.waitGroup.Wait()
}

func Go(goroutine func()) {
	GoWithRecover(goroutine, defaultRecoverGoroutine)
}

// GoWithRecover starts a recoverable goroutine using given customRecover() function.
func GoWithRecover(goroutine func(), customRecover func(err interface{})) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}

func defaultRecoverGoroutine(err interface{}) {
	logger.Errorf("Error in Go routine: %v", err)
	logger.Errorf("Stack: %s", debug.Stack())
}
