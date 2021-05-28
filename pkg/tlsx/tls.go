// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/5/17

package tlsx

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"github.com/crochee/proxy-go/pkg/filecontent"
)

// TlsConfig output tls
func TlsConfig(clientAuth tls.ClientAuthType, ca, cert, key filecontent.FileOrContent) (*tls.Config, error) {
	caPEMBlock, err := ca.Read()
	if err != nil {
		return nil, err
	}
	var certPEMBlock []byte
	if certPEMBlock, err = cert.Read(); err != nil {
		return nil, err
	}
	var keyPEMBlock []byte
	if keyPEMBlock, err = key.Read(); err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEMBlock) {
		return nil, errors.New("failed to parse root certificate")
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   clientAuth,
		ClientCAs:    pool,
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
