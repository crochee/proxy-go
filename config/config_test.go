// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"os"
	"proxy-go/config/dynamic"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestInitConfig(t *testing.T) {
	cf := &Config{
		Server: &Server{
			Medata: []*Medata{
				{
					GraceTimeOut: 10 * time.Second,
					Name:         "server1",
					Scheme:       "http",
					Port:         8120,
				},
				{
					GraceTimeOut: 10 * time.Second,
					Name:         "server2",
					Scheme:       "https",
					Port:         8121,
					Tls: &TlsConfig{
						Cert: "./conf/cert.pem",
						Key:  "./conf/key.pem",
					},
				},
			},
		},
		Middleware: &dynamic.Middleware{
			ReplaceHost: &dynamic.ReplaceHost{
				Scheme: "http",
				Host:   "127.0.0.1:8150",
			},
		},
	}

	file, err := os.Create("../conf/config.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err = yaml.NewEncoder(file).Encode(cf); err != nil {
		t.Fatal(err)
	}
}
