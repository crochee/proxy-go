// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package internal

import (
	"net"
	"os"
	"testing"
)

func TestGenerateSelfSignedCertKey(t *testing.T) {
	//randomBytes := make([]byte, 100)
	//if _, err := rand.Read(randomBytes); err != nil {
	//	t.Error(err)
	//}
	//zBytes := sha256.Sum256(randomBytes)
	//z := hex.EncodeToString(zBytes[:sha256.Size])
	//domain := fmt.Sprintf("%s.%s.proxy.default", z[:32], z[32:])
	//domain := "localhost"
	host := "192.168.31.62"
	cert, key, err := GenerateSelfSignedCertKey(
		host,
		[]net.IP{
			net.ParseIP(host),
		},
		[]string{
			//	domain,
		})
	if err != nil {
		t.Error(err)
	}
	certFile, err := os.Create("../conf/cert.pem")
	if err != nil {
		t.Error(err)
	}
	if _, err = certFile.Write(cert); err != nil {
		t.Error(err)
	}
	if err = certFile.Close(); err != nil {
		t.Error(err)
	}
	keyFile, err := os.Create("../conf/key.pem")
	if err != nil {
		t.Error(err)
	}
	if _, err = keyFile.Write(key); err != nil {
		t.Error(err)
	}
	if err = keyFile.Close(); err != nil {
		t.Error(err)
	}

}
