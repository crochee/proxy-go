// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package cmd

import (
	"context"
	"crypto/tls"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/router"
	"github.com/crochee/proxy-go/pkg/tlsx"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
	"github.com/crochee/proxy-go/pkg/transport/pprofx"
	"github.com/crochee/proxy-go/pkg/transport/prometheusx"
)

func Server(cmd *cobra.Command, _ []string) error {
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
	var tlsConfig *tls.Config
	if tlsConfig, err = tlsx.TlsConfig(tls.NoClientCert, config.Cfg.Medata.Tls.Ca,
		config.Cfg.Medata.Tls.Cert, config.Cfg.Medata.Tls.Key); err != nil {
		return err
	}
	var serverList []transport.AppServer
	httpSrv, err := httpx.New(ctx, config.Cfg.Medata.Host, router.Handler(config.Cfg),
		httpx.TlsConfig(tlsConfig), httpx.RequestLog(logger.NewLogger(logger.Path(config.Cfg.Medata.RequestLog.Path),
			logger.Level(config.Cfg.Medata.RequestLog.Level))))
	if err != nil {
		return err
	}
	serverList = append(serverList, httpSrv)
	if config.Cfg.PrometheusAgent != "" {
		serverList = append(serverList, prometheusx.New(ctx, config.Cfg.PrometheusAgent))
	}
	if config.Cfg.PprofAgent != "" {
		serverList = append(serverList, pprofx.New(ctx, config.Cfg.PprofAgent))
	}
	app := transport.NewApp(
		transport.Context(ctx),
		transport.Servers(serverList...),
	)
	if err = app.Run(); err != nil {
		return err
	}
	for _, srv := range serverList {
		logger.Infof("server %s stop", srv.Name())
	}
	logger.Exit("proxy exit!")
	return nil
}
