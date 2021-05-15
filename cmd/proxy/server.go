// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
	"github.com/crochee/proxy-go/router"
)

func server(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel() // 全局取消
	// 初始化配置
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}
	config.InitConfig(configFile)
	// 初始化系统日志
	if config.Cfg.Medata.SystemLog != nil {
		logger.InitSystemLogger(logger.Path(config.Cfg.Medata.SystemLog.Path),
			logger.Level(config.Cfg.Medata.SystemLog.Level))
	}

	httpSrv, err := httpx.New(ctx, config.Cfg.Medata, router.Handler(config.Cfg))
	if err != nil {
		return err
	}
	app := transport.NewApp(
		transport.Context(ctx),
		transport.Servers(httpSrv),
	)
	if err = app.Run(); err != nil {
		return err
	}
	logger.Exit("server exit!")
	return nil
}
