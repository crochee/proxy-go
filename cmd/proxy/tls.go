// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/18

package main

import (
	"net"
	"os"

	"github.com/urfave/cli/v2"

	"proxy-go/ptls"
)

func certificate(c *cli.Context) error {
	host := c.String("host")
	domain := c.String("domain")
	cert, key, err := ptls.GenerateSelfSignedCertKey(
		host,
		[]net.IP{
			net.ParseIP(host),
		},
		[]string{
			domain,
		})
	if err != nil {
		return err
	}
	var certFile *os.File
	if certFile, err = os.Create(c.String("cert")); err != nil {
		return nil
	}
	defer certFile.Close()
	if _, err = certFile.Write(cert); err != nil {
		return err
	}
	var keyFile *os.File
	if keyFile, err = os.Create(c.String("key")); err != nil {
		return err
	}
	defer keyFile.Close()
	_, err = keyFile.Write(key)
	return err
}
