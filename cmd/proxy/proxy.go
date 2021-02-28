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

	"github.com/crochee/proxy-go/config"
	"github.com/crochee/proxy-go/logger"
	"github.com/crochee/proxy-go/router"
	"github.com/crochee/proxy-go/safe"
	"github.com/crochee/proxy-go/server"
	"github.com/crochee/proxy-go/version"
)

func main() {
	app := cli.NewApp()
	app.Name = "proxy"
	app.Version = version.Version
	app.Usage = "Generates proxy"

	app.Commands = []*cli.Command{
		{
			Name:    "proxy",
			Aliases: []string{"p"},
			Usage:   "proxy server",
			Action:  run,
			Flags: []cli.Flag{
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
			},
		},
		{
			Name:    "tls",
			Aliases: []string{"t"},
			Usage:   "generates random TLS certificates",
			Action:  certificate,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "cert",
					Aliases: []string{"c"},
					Usage:   "cert path",
					EnvVars: []string{"cert_path"},
					Value:   "./conf/cert.pem",
				},
				&cli.StringFlag{
					Name:    "key",
					Aliases: []string{"k"},
					Usage:   "key path",
					EnvVars: []string{"key_path"},
					Value:   "./conf/key.pem",
				},
				&cli.StringFlag{
					Name:    "host",
					Aliases: []string{"h"},
					Usage:   "host",
					EnvVars: []string{"host"},
					Value:   "127.0.0.1",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "domain",
					EnvVars: []string{"domain"},
					Value:   "localhost",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
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

	handler, err := router.ChainBuilder(server.GlobalWatcher)
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
