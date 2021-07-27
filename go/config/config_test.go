package config

import (
	"testing"
	"time"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/pkg/tlsx"
)

func TestInitConfig(t *testing.T) {
	cf := &Spec{
		Medata: &Medata{
			Tls: &tlsx.Config{
				Ca:   "./build/package/proxy/cert/ca.pem",
				Cert: "./build/package/proxy/cert/proxy.pem",
				Key:  "./build/package/proxy/cert/proxy-key.pem",
			},
			GraceTimeOut: 15 * time.Second,
			Scheme:       "https",
			Host:         ":8120",
			SystemLog: &dynamic.LogInfo{
				Path:  "./log/sys.log",
				Level: "DEBUG",
			},
		},
		Server: &Medata{
			Tls: &tlsx.Config{
				Ca:   "./build/package/proxy/cert/ca.pem",
				Cert: "./build/package/proxy/cert/server.pem",
				Key:  "./build/package/proxy/cert/server-key.pem",
			},
			GraceTimeOut: 5 * time.Second,
			Scheme:       "https",
			Host:         ":8121",
			RequestLog: &dynamic.LogInfo{
				Path:  "./log/proxy.log",
				Level: "INFO",
			},
		},
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
