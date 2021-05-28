// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/28

// Package main
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/vugu/vugu/simplehttp"

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
	"github.com/crochee/proxy-go/pkg/transport/prometheusx"
)

var (
	host = flag.String("host", "localhost:8121", "web")
	mode = flag.Bool("mode", true, "web dev mode")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消

	httpSrv, err := httpx.New(ctx, *host, simplehttp.New(*host, *mode))
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app := transport.NewApp(
		transport.Context(ctx),
		transport.Servers(
			httpSrv,
			prometheusx.New(ctx, config.Cfg.PrometheusAgent.Host, config.Cfg.PrometheusAgent.Path),
		),
	)
	if err = app.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	logger.Exit("server exit!")
}
