// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"proxy-go/cmd"
	"proxy-go/config"
	"proxy-go/logger"
	"proxy-go/router"
	"proxy-go/safe"
	"proxy-go/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "proxy"
	app.Version = cmd.Version
	app.Usage = "Generates proxy"

	app.Commands = cli.Commands{
		{
			Name:    "proxy",
			Aliases: []string{"p"},
			Usage:   "proxy server",
			Action:  run,
			Flags:   runFlags,
		},
		{
			Name:    "tls",
			Aliases: []string{"t"},
			Usage:   "generates random TLS certificates",
			Action:  certificate,
			Flags:   TlsFlags,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var runFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "enable-log",
		Usage:   "enable log switch",
		EnvVars: []string{"enable_log"},
	},
	&cli.StringFlag{
		Name:    "log-path",
		Usage:   "log path",
		EnvVars: []string{"log_path"},
	},
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"log_level"},
	},
	&cli.StringFlag{
		Name:    "config",
		Usage:   "config path",
		EnvVars: []string{"config"},
		Value:   "./conf/config.yml",
	},
}

func run(c *cli.Context) error {
	ctx := logger.With(context.Background(),
		logger.Enable(c.Bool("enable-log")),
		logger.Level(strings.ToUpper(c.String("log-level"))),
		logger.LogPath(c.String("log-path")),
	)
	path := c.String("config")
	if path == "" {
		path = "./conf/config.yml"
	}
	config.InitConfig(path)
	return setup(ctx)
}

func setup(ctx context.Context) error {
	ctx = server.ContextWithSignal(ctx)
	pool := safe.NewPool(ctx)

	server.GlobalWatcher = server.NewWatcher(ctx, pool)

	handler, err := router.ChainBuilder(ctx, server.GlobalWatcher)
	if err != nil {
		logger.FromContext(ctx).Fatalf("build route failed.Error:%v", err)
		return err
	}
	srv := server.NewServer(ctx, config.Cfg, handler, server.GlobalWatcher)

	srv.Start()
	defer srv.Close()

	srv.Wait()
	logger.FromContext(ctx).Info("shutting down")
	return nil
}
