package main

import (
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/crochee/proxy-go/internal"
)

func TlsTool(cmd *cobra.Command, _ []string) error {
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
	if cert, key, err = internal.GenerateSelfSignedCertKey(
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
		return err
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
