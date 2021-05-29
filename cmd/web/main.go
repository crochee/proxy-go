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

	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/pkg/transport"
	"github.com/crochee/proxy-go/pkg/transport/httpx"
)

var (
	host = flag.String("host", "localhost:8121", "web")
	mode = flag.Bool("mode", true, "web dev mode")
)

func main() {
	if err := mainFunc(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func mainFunc() error {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 全局取消

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	httpSrv, err := httpx.New(ctx, *host, simplehttp.New(wd, *mode))
	if err != nil {
		return err
	}
	app := transport.NewApp(
		transport.Context(ctx),
		transport.Servers(
			httpSrv,
		),
	)
	if err = app.Run(); err != nil {
		return err
	}
	logger.Exit("server exit!")
	return nil
}
