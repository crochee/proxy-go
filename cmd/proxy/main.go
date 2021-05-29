// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/10

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/cmd"
	"github.com/crochee/proxy-go/version"
)

func main() {
	rootCmd := &cobra.Command{
		Short:   "proxy tools",
		Version: version.Version,
	}

	serverCmd := &cobra.Command{
		Use:    "server",
		Short:  "start server",
		Long:   "start multi server",
		RunE:   cmd.Server,
		Hidden: true,
	}
	serverCmd.Flags().StringP("config", "c", "./conf/config.yml", "")

	tlsCmd := &cobra.Command{
		Use:   "tls",
		Short: "generate tls file",
		Long:  "generate self tls file",
		RunE:  cmd.TlsTool,
	}
	tlsCmd.Flags().StringP("ip", "i", "127.0.0.1", "")
	tlsCmd.Flags().StringP("domain", "d", "localhost", "")
	tlsCmd.Flags().StringP("cert", "c", "./conf/cert.pem", "")
	tlsCmd.Flags().StringP("key", "k", "./conf/key.pem", "")

	rootCmd.AddCommand(tlsCmd, serverCmd)
	rootCmd.SetErr(bytes.NewBuffer(nil))

	if err := rootCmd.Execute(); err != nil && !errors.Is(err, context.Canceled) {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
