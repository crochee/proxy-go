// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/16

package routine

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
)

type pool struct {
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	option
}

// NewPool creates a Pool.
func NewPool(parentCtx context.Context, opts ...func(*option)) *pool {
	ctx, cancel := context.WithCancel(parentCtx)
	p := &pool{
		ctx:    ctx,
		cancel: cancel,
		option: option{recoverFunc: defaultRecoverGoroutine},
	}
	for _, opt := range opts {
		opt(&p.option)
	}
	return p
}

// Go starts a recoverable goroutine with a context.
func (p *pool) Go(goroutine func(context.Context)) {
	p.waitGroup.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if p.recoverFunc != nil {
					p.recoverFunc(r)
				}
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

func defaultRecoverGoroutine(err interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "Error:%v\nStack: %s", err, debug.Stack())
}
