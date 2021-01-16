// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"proxy-go/safe"
	"strings"
	"sync"
	"syscall"
	"time"

	"proxy-go/config"
	"proxy-go/logger"
)

func main() {
	initContext()

	ctx := context.Background()
	pool := safe.NewPool(ContextWithSignal(ctx))

	if config.Cfg.Server.Port == nil {
		logger.Fatal("no config port to start")
	}
	var handler http.Handler = ChanCreate()
	// http
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Cfg.Server.Port.Https),
		Handler: handler,
	}
	if config.Cfg.Server.Port.Http != 0 {
		pool.GoCtx(func(ctx context.Context) {
			logger.Info("http server running...")
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error(err.Error())
			}
		})
	}
	httpsSrv := &http.Server{
		Addr:      fmt.Sprintf(":%d", config.Cfg.Server.Port.Https),
		Handler:   handler,
		TLSConfig: nil,
	}
	// https
	if config.Cfg.Server.Port.Https != 0 {
		go func() {
			logger.Info("https server running...")
			if err := httpsSrv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
				logger.Fatal(err.Error())
			}
		}()
	}
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
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		Shutdown(ctx, httpSrv)
		wg.Done()
	}()
	go func() {
		Shutdown(ctx, httpsSrv)
		wg.Done()
	}()
	wg.Wait()

	pool.Stop()
	logger.Info("server exit!")
}

func initContext() {
	// 配置路径获取
	path, ok := os.LookupEnv("config")
	if !ok {
		path = "./conf/config.yml"
	}
	config.InitConfig(path)

	// 日志初始化
	logger.InitLogger(
		logger.Enable(os.Getenv("enable-log") == "true"),
		logger.Level(strings.ToUpper(os.Getenv("log-level"))),
		logger.LogPath(os.Getenv("log-path")),
	)
}

// a channel (just for the fun of it)
type Chan chan int

func ChanCreate() Chan {
	c := make(Chan)
	go func(c Chan) {
		for x := 0; ; x++ {
			c <- x
		}
	}(c)
	return c
}

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fmt.Sprintf("channel send #%d\n", <-ch))
}

func Shutdown(ctx context.Context, server *http.Server) {
	err := server.Shutdown(ctx)
	if err == nil {
		return
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logger.Debugf("server failed to shutdown within deadline because: %s", err)
		if err = server.Close(); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	logger.Error(err.Error())
	// We expect Close to fail again because Shutdown most likely failed when trying to close a listener.
	// We still call it however, to make sure that all connections get closed as well.
	server.Close()
}
