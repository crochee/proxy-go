// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/16

package routine

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
)

type errGroup struct {
	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
	errOnce   sync.Once
	err       error
}

// NewGroup starts a recoverable goroutine errGroup with a context.
func NewGroup(ctx context.Context) *errGroup {
	newCtx, cancel := context.WithCancel(ctx)
	return &errGroup{
		ctx:    newCtx,
		cancel: cancel,
	}
}

// Go starts a recoverable goroutine with a context.
func (e *errGroup) Go(goroutine func(context.Context) error) {
	e.waitGroup.Add(1)
	go func() {
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v.Stack:%s", r, debug.Stack())
			}
			if err != nil {
				e.errOnce.Do(func() {
					e.err = err
					e.cancel()
				})
			}
			e.waitGroup.Done()
		}()
		err = goroutine(e.ctx)
	}()
}

func (e *errGroup) Wait() error {
	e.waitGroup.Wait()
	e.errOnce.Do(e.cancel)
	return e.err
}
