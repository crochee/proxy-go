// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// ContextWithSignal creates a context canceled when SIGINT or SIGTERM are notified.
func ContextWithSignal(ctx context.Context) context.Context {
	newCtx, cancel := context.WithCancel(ctx)
	signals := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(signals, syscall.SIGINT)
	go func() {
		<-signals
		cancel()
	}()
	return newCtx
}
