// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/4

package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/cmd"
	"github.com/crochee/proxy-go/ptls"
)

func main() {
	rootCmd := &cobra.Command{
		Args:    cobra.MinimumNArgs(1),
		Version: cmd.Version,
		RunE:    run,
	}
	rootCmd.Flags().StringP("ip", "i", "127.0.0.1", "")
	rootCmd.Flags().StringP("domain", "d", "localhost", "")
	rootCmd.Flags().StringP("cert", "c", "./conf/cert.pem", "")
	rootCmd.Flags().StringP("key", "k", "./conf/key.pem", "")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func run(cmd *cobra.Command, args []string) error {
	flag := cmd.Flags()
	host, err := flag.GetString("ip")
	if err != nil {
		return err
	}
	var domain string
	if domain, err = flag.GetString("domain"); err != nil {
		return err
	}
	var (
		cert []byte
		key  []byte
	)
	if cert, key, err = ptls.GenerateSelfSignedCertKey(
		host,
		[]net.IP{
			net.ParseIP(host),
		},
		[]string{
			domain,
		}); err != nil {
		return err
	}
	var certPath string
	if certPath, err = flag.GetString("cert"); err != nil {
		return err
	}
	var keyPath string
	if keyPath, err = flag.GetString("key"); err != nil {
		return err
	}
	var certFile *os.File
	if certFile, err = os.Create(certPath); err != nil {
		return nil
	}
	defer certFile.Close()
	if _, err = certFile.Write(cert); err != nil {
		return err
	}
	var keyFile *os.File
	if keyFile, err = os.Create(keyPath); err != nil {
		return err
	}
	defer keyFile.Close()
	_, err = keyFile.Write(key)
	return err
}
