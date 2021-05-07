// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"context"
	"flag"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/router"
	"github.com/crochee/proxy-go/service/transport"
	"github.com/crochee/proxy-go/service/transport/httpx"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(*configFile)
	// 初始化系统日志
	if config.Cfg.Medata.SystemLog != nil {
		logger.InitSystemLogger(logger.Path(config.Cfg.Medata.SystemLog.Path),
			logger.Level(config.Cfg.Medata.SystemLog.Level))
	}

	httpSrv, err := httpx.New(ctx, config.Cfg.Medata, router.Handler(config.Cfg.Middleware))
	if err != nil {
		logger.Fatal(err.Error())
	}
	app := transport.NewApp(
		transport.Context(ctx),
		transport.Servers(httpSrv),
	)
	if err := app.Run(); err != nil {
		logger.Error(err.Error())
	}
	logger.Exit("server exit!")
}
