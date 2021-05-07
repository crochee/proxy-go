// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"testing"
	"time"

	"github.com/crochee/proxy-go/config/dynamic"
)

func TestInitConfig(t *testing.T) {
	cf := &Spec{
		Medata: &Medata{
			Tls:          nil,
			GraceTimeOut: 15 * time.Second,
			Scheme:       "http",
			Host:         ":8120",
			SystemLog: &dynamic.LogInfo{
				Path:  "./log/proxy-sys.log",
				Level: "DEBUG",
			},
			RequestLog: &dynamic.LogInfo{
				Path:  "./log/proxy-req.log",
				Level: "DEBUG",
			},
		},
		Middleware: &dynamic.Config{
			Balance: map[string]*dynamic.Balance{
				"obs": {
					Selector: "wwr",
					NodeList: []*dynamic.Node{
						{
							Scheme:   "http",
							Host:     "127.0.0.1:8121",
							Metadata: nil,
							Weight:   1.0,
						},
						{
							Scheme:   "http",
							Host:     "127.0.0.1:8122",
							Metadata: nil,
							Weight:   2.0,
						},
					},
				},
			},
			AccessLog: &dynamic.LogInfo{
				Path:  "",
				Level: "INFO",
			},
			RateLimit: &dynamic.RateLimit{
				Every: time.Second,
				Burst: 2000,
				Mode:  1,
			},
			Recovery:    true,
			CrossDomain: false,
		},
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
