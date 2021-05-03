// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"net"
	"net/http"
	"strings"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
)

var configFile = flag.String("f", "./conf/config.yml", "the config file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消
	// 初始化配置
	config.InitConfig(*configFile)
	// 初始化系统日志
	pathFunc := func(option *logger.Option) {
		if config.Cfg.Medata.LogPath == "" {
			return
		}
		option.Path = config.Cfg.Medata.LogPath
	}
	levelFunc := func(option *logger.Option) {
		if config.Cfg.Medata.LogLevel == "" {
			return
		}
		option.Level = config.Cfg.Medata.LogLevel
	}
	logger.InitSystemLogger(pathFunc, levelFunc)
	// 初始化请求日志
	requestLog := logger.NewLogger(pathFunc, levelFunc)

	ln, err := net.Listen("tcp", config.Cfg.Medata.Host)
	if err != nil {
		logger.Fatal(err.Error())
	}
	srv := &http.Server{
		Handler: nil,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return logger.With(ctx, requestLog)
		},
	}
	logger.Infof("server medata:%+v running...", config.Cfg.Medata)
	switch strings.ToLower(config.Cfg.Medata.Scheme) {
	case "http":
	case "https":
		if ln, err = tlsListener(ln, config.Cfg.Medata); err != nil {
			logger.Fatal(err.Error())
		}
	default:
		logger.Fatalf("scheme is %s", config.Cfg.Medata.Scheme)
	}
	if err = srv.Serve(ln); err != nil {
		logger.Fatal(err.Error())
	}
}

func tlsListener(listener net.Listener, medata *config.Medata) (net.Listener, error) {
	if medata.Tls == nil {
		return nil, errors.New("https haven't tls")
	}
	certPEMBlock, err := medata.Tls.Cert.Read()
	if err != nil {
		return nil, err
	}
	var keyPEMBlock []byte
	if keyPEMBlock, err = medata.Tls.Key.Read(); err != nil {
		return nil, err
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
		return nil, err
	}

	return tls.NewListener(listener, &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}), nil
}
