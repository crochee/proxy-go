package config

import (
	"testing"
	"time"

	"github.com/crochee/proxy-go/config/dynamic"
)

func TestInitConfig(t *testing.T) {
	cf := &Spec{
		Medata: &Medata{
			Tls: &TlsConfig{
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
			RequestLog: &dynamic.LogInfo{
				Path:  "./log/req.log",
				Level: "DEBUG",
			},
		},
		Middleware: &dynamic.Config{
			Retry: &dynamic.Retry{
				Attempts:        10,
				InitialInterval: 5 * time.Second,
			},
			AccessLog: &dynamic.LogInfo{
				Path:  "./log/access.log",
				Level: "INFO",
			},
			Trace: nil,
			Balance: &dynamic.BalanceConfig{
				RegisterApis: []*dynamic.ServiceApi{
					{
						ServiceName: "OBS",
						Path:        "/proxy",
						Method:      "POST",
					},
				},
				Transfers: []*dynamic.ServiceBalance{
					{
						ServiceName: "OBS",
						Balance: dynamic.Balance{
							Selector: "wrr",
							Nodes: []*dynamic.Node{
								{
									Scheme:   "https",
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
				},
			},
			RateLimit: &dynamic.RateLimit{
				Every: time.Second,
				Burst: 2000,
				Mode:  1,
			},
			CrossDomain:    false,
			CircuitBreaker: &dynamic.CircuitBreaker{Expression: "NetworkErrorRatio() > 0.5"},
			Recovery:       true,
		},
		Proxy: &Proxy{
			ProxyLog: &dynamic.LogInfo{
				Path:  "./log/proxy.log",
				Level: "INFO",
			},
			Tls: &TlsConfig{
				Ca:   "./build/package/proxy/cert/ca.pem",
				Cert: "./build/package/proxy/cert/proxy.pem",
				Key:  "./build/package/proxy/cert/proxy-key.pem",
			},
		},
		PrometheusAgent: ":8190",
		PprofAgent:      ":8191",
	}
	y := Yml{path: "../conf/config.yml"}
	if err := y.Encode(cf); err != nil {
		t.Error(err)
	}
}
