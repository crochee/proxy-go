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
			GraceTimeOut: 0,
			Scheme:       "",
			Host:         ":8120",
		},
		Middleware: &dynamic.Config{
			Balance: map[string]*dynamic.Balance{
				"obs": {
					Selector: "wwr",
					NodeList: []*dynamic.Node{
						{
							Scheme:   "http",
							Host:     "127.0.0.1:8701",
							Metadata: nil,
							Weight:   0,
						},
					},
				},
			},
			RateLimit: &dynamic.RateLimit{
				Every: 10 * time.Second,
				Burst: 2000,
				Mode:  0,
			},
		},
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
