// Copyright (c) Huawei Technologies Co., Ltd. 2021-2021. All rights reserved.
// Description:
// Author: licongfu
// Create: 2021/5/10

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/cmd"
)

func main() {
	rootCmd := &cobra.Command{
		Short:   "proxy tools",
		Version: cmd.Version,
	}
	tlsCmd := &cobra.Command{
		Use:   "tls",
		Short: "generate tls file",
		Long:  "generate self tls file",
		RunE:  tls,
	}
	tlsCmd.Flags().StringP("ip", "i", "127.0.0.1", "")
	tlsCmd.Flags().StringP("domain", "d", "localhost", "")
	tlsCmd.Flags().StringP("cert", "c", "./conf/cert.pem", "")
	tlsCmd.Flags().StringP("key", "k", "./conf/key.pem", "")

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "start server",
		Long:  "start multi server",
		RunE:  server,
	}

	serverCmd.Flags().StringP("config", "c", "./conf/config.yml", "")

	rootCmd.AddCommand(tlsCmd, serverCmd)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
