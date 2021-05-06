// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/crochee/proxy-go/safe"
)

type Server interface {
	Start() error
	Stop() error
}

type App struct {
	option
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp(opts ...func(*option)) *App {
	option := option{
		sigList: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		ctx:     context.Background(),
	}
	for _, opt := range opts {
		opt(&option)
	}
	ctx, cancel := context.WithCancel(option.ctx)
	return &App{
		option: option,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() error {
	g := safe.NewErrGroup(a.ctx)
	for _, srv := range a.serverList {
		realSrv := srv
		g.Go(func(ctx context.Context) error {
			<-ctx.Done()
			return realSrv.Stop()
		})
		g.Go(func(ctx context.Context) error {
			return realSrv.Start()
		})
	}
	g.Go(func(ctx context.Context) error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, a.option.sigList...)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				if a.cancel != nil {
					a.cancel()
				}
			}
		}
	})
	return g.Wait()
}

func (a *App) stop() {
	if a.cancel != nil {
		a.cancel()
	}
}
