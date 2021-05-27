// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/6

// Package safe
package safe

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/crochee/proxy-go/logger"
)

type errGroup struct {
	pool    *pool
	errOnce sync.Once
	err     error
}

func NewErrGroup(ctx context.Context) *errGroup {
	return &errGroup{
		pool: NewPool(ctx),
	}
}

// Go starts a recoverable goroutine with a context.
func (e *errGroup) Go(goroutine func(ctx context.Context) error) {
	e.pool.waitGroup.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.FromContext(e.pool.ctx).Errorf("[Recovery] panic happened.Error:%v\n.Stack:\n%s", debug.Stack())
			}
			e.pool.waitGroup.Done()
		}()
		if err := goroutine(e.pool.ctx); err != nil {
			e.errOnce.Do(func() {
				e.err = err
				if e.pool.cancel != nil {
					e.pool.cancel()
				}
			})
		}
	}()
}

func (e *errGroup) Wait() error {
	e.pool.waitGroup.Wait()
	e.errOnce.Do(func() {
		if e.pool.cancel != nil {
			e.pool.cancel()
		}
	})
	return e.err
}
