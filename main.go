// Copyright 2021, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package main

import (
	"context"
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

	ctx = server.ContextWithSignal(ctx)
	pool := safe.NewPool(ctx)

	handler, err := router.Route(ctx, config.Cfg)
	if err != nil {
		logger.FromContext(ctx).Fatalf("build route failed.Error:%v", err)
	}
	srv := server.NewServer(ctx, pool, config.Cfg, handler)

	srv.Start()
	defer srv.Close()

	srv.Wait()
	logger.FromContext(ctx).Info("shutting down")
}

func initContext() {
	// 配置路径获取
	path, ok := os.LookupEnv("config")
	if !ok {
		path = "./conf/config.yml"
	}
	config.InitConfig(path)
}
