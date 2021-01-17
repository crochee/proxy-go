// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"proxy-go/config"
	"proxy-go/logger"
	"proxy-go/router"
	"proxy-go/safe"
	"proxy-go/server"
)

func main() {
	initContext()
	// 日志初始化
	ctx := logger.With(
		context.Background(),
		logger.Enable(os.Getenv("enable_log") == "true"),
		logger.Level(strings.ToUpper(os.Getenv("log_level"))),
		logger.LogPath(os.Getenv("log_path")),
	)

	ctx = ContextWithSignal(ctx)
	pool := safe.NewPool(ctx)

	handler, err := router.Route(ctx)
	if err != nil {
		logger.FromContext(ctx).Fatalf("build route failed.Error:%v", err)
	}
	srv := server.New(ctx, pool, config.Cfg, handler)

	srv.Start()
	defer srv.Close()

	srv.Wait()
	logger.FromContext(ctx).Info("Shutting down")
}

func initContext() {
	// 配置路径获取
	path, ok := os.LookupEnv("config")
	if !ok {
		path = "./conf/config.yml"
	}
	config.InitConfig(path)

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
