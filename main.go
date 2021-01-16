// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxy-go/config"
	"proxy-go/logger"
)

func main() {
	config.InitConfig()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 1),
		Handler: nil,
	}
	go func() {
		logger.Info("server running...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("server forced to shutdown:%v", err)
	}
	logger.Info("server exit!")
}
