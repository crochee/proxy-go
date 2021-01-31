// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"testing"
	"time"
)

func TestInitConfig(t *testing.T) {
	cf := &Config{
		Server: &Server{
			Medata: []*Medata{
				{
					Name:         "proxy1",
					Scheme:       "http",
					Port:         8120,
					Tls:          nil,
					GraceTimeOut: 15 * time.Second,
				},
				{
					Name:   "proxy2",
					Scheme: "https",
					Port:   8121,
					Tls: &TlsConfig{
						Cert: "./conf/cert.pem",
						Key:  "./conf/key.pem",
					},
					GraceTimeOut: 15 * time.Second,
				},
			},
		},
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
