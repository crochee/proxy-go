// Package proxygo
package proxygo

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/pkg/logger"
	"github.com/crochee/proxy-go/pkg/metrics"
	"github.com/crochee/proxy-go/pkg/middleware"
	"github.com/crochee/proxy-go/pkg/router"
	"github.com/crochee/proxy-go/pkg/tlsx"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
)

func Server() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消

	// 初始化系统日志
	if config.Cfg.Medata.SystemLog != nil {
		logger.InitSystemLogger(logger.Path(config.Cfg.Medata.SystemLog.Path),
			logger.Level(config.Cfg.Medata.SystemLog.Level))
	}

	var serverList []transport.AppServer
	proxyHttp, err := httpAppServer(ctx, config.Cfg.Medata, middleware.Load(ctx))
	if err != nil {
		return err
	}
	serverList = append(serverList, proxyHttp)

	var opts []httpx.Option
	opts = append(opts, httpx.BeforeStart(func() error {
		// 注册普罗米修斯
		metrics.DefineMetrics()
		prometheus.MustRegister(metrics.ReqDurHistogramVec, metrics.ReqCodeTotalCounter)
		return nil
	}))
	opts = append(opts, httpx.AfterStop(func() error {
		// 卸载普罗米修斯
		prometheus.Unregister(metrics.ReqDurHistogramVec)
		prometheus.Unregister(metrics.ReqCodeTotalCounter)
		return nil
	}))
	var serverHttp transport.AppServer
	if serverHttp, err = httpAppServer(ctx, config.Cfg.Server, router.ApiHandler(), opts...); err != nil {
		return err
	}
	serverList = append(serverList, serverHttp)

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
	logger.Exit("system exit!")
	return nil
}

func httpAppServer(ctx context.Context, medata *config.Medata,
	handler http.Handler, opts ...httpx.Option) (transport.AppServer, error) {
	if medata == nil {
		return nil, errors.New("medata is nil")
	}
	if medata.Tls != nil {
		tlsConfig, err := tlsx.TlsConfig(tls.NoClientCert, *medata.Tls)
		if err != nil {
			return nil, err
		}
		opts = append(opts, httpx.TlsConfig(tlsConfig))
	}
	if medata.RequestLog != nil {
		opts = append(opts, httpx.RequestLog(logger.NewLogger(logger.Path(medata.RequestLog.Path),
			logger.Level(medata.RequestLog.Level))))
	}
	return httpx.New(ctx, medata.Host, handler, opts...)
}
