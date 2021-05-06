// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"context"
	"flag"
	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/service/server"
	"github.com/crochee/proxy-go/service/server/httpx"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(*configFile)
	// 初始化系统日志
	logger.InitSystemLogger(logger.Path(config.Cfg.Medata.LogPath), logger.Level(config.Cfg.Medata.LogLevel))

	httpSrv, err := httpx.New(ctx, config.Cfg.Medata, nil)
	if err != nil {
		logger.Fatal(err.Error())
	}
	app := server.NewApp(
		server.Context(ctx),
		server.Servers(httpSrv),
	)
	if err := app.Run(); err != nil {
		logger.Error(err.Error())
	}
	logger.Exit("server exit!")
}
