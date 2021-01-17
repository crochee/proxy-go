// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package ptls

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

// DefaultDomain proxy domain for the default certificate.
const DefaultDomain = "PROXY DEFAULT CERT"

// DefaultCertificate generates random TLS certificates.
func DefaultCertificate() (*tls.Certificate, error) {
	certPEM, keyPEM, err := GenerateSelfSignedCertKey("127.0.0.1", nil, []string{DefaultDomain})
	if err != nil {
		return nil, err
	}
	var certificate tls.Certificate
	if certificate, err = tls.X509KeyPair(certPEM, keyPEM); err != nil {
		return nil, err
	}
	return &certificate, nil
}

func GenerateSelfSignedCertKey(host string, alternateIPs []net.IP, alternateDNS []string) ([]byte, []byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}
	var privateKey *rsa.PrivateKey
	if privateKey, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
		return nil, nil, err
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: fmt.Sprintf("%s@%d", host, now.Unix()),
		},
		NotBefore: now,
		NotAfter:  now.Add(time.Hour * 24 * 365),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	template.IPAddresses = append(template.IPAddresses, alternateIPs...)
	template.DNSNames = append(template.DNSNames, alternateDNS...)

	var derBytes []byte
	if derBytes, err = x509.CreateCertificate(rand.Reader,
		&template,
		&template,
		&privateKey.PublicKey,
		privateKey); err != nil {
		return nil, nil, err
	}

	// Generate cert
	var certBuffer bytes.Buffer
	if err := pem.Encode(&certBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}); err != nil {
		return nil, nil, err
	}

	// Generate key
	var keyBuffer bytes.Buffer
	if err := pem.Encode(&keyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return nil, nil, err
	}

	return certBuffer.Bytes(), keyBuffer.Bytes(), nil
}
