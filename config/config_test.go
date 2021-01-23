// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"testing"
	"time"

	"proxy-go/config/dynamic"
)

func TestInitConfig(t *testing.T) {
	cf := &Config{
		Server: &Server{
			Medata: []*dynamic.Medata{
				{
					Name:         "proxy",
					Scheme:       "http",
					Port:         8120,
					Tls:          nil,
					GraceTimeOut: 15 * time.Second,
					Mode:         "Random",
					Path:         "/obs/",
					LocationList: []*dynamic.Location{
						{
							ProxyPass: "http://127.0.0.1:8150/v1/",
							Weight:    1,
						},
					},
					Middleware: nil,
				},
			},
		},
	}
	y := Json{path: "../conf/config.json"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
