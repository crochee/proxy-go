package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"os"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/router"
	"github.com/crochee/proxy-go/pkg/tlsx"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
	"github.com/crochee/proxy-go/pkg/transport/pprofx"
	"github.com/crochee/proxy-go/pkg/transport/prometheusx"
)

var configFile = flag.String("config", "./conf/config.yml", "")

func main() {
	if err := server(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func server() error {
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

	tlsConfig, err := tlsx.TlsConfig(tls.NoClientCert, config.Cfg.Medata.Tls.Ca,
		config.Cfg.Medata.Tls.Cert, config.Cfg.Medata.Tls.Key)
	if err != nil {
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
