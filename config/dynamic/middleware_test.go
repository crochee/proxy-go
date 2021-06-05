package dynamic

import (
	"os"
	"testing"
	"time"

	"github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/crochee/proxy-go/pkg/tlsx"
)

func TestNew(t *testing.T) {
	cf := &Config{
		Middleware: &Middleware{
			Retry: &Retry{
				Attempts:        10,
				InitialInterval: 5 * time.Second,
			},
			AccessLog: &LogInfo{
				Path:  "./log/access.log",
				Level: "INFO",
			},
			Trace: nil,
			Balance: &BalanceConfig{
				RegisterApis: []*ServiceApi{
					{
						ServiceName: "OBS",
						Path:        "/proxy",
						Method:      "POST",
					},
				},
				Transfers: []*ServiceBalance{
					{
						ServiceName: "OBS",
						Balance: Balance{
							Selector: "wrr",
							Nodes: []*Node{
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
			RateLimit: &RateLimit{
				Every: time.Second,
				Burst: 2000,
				Mode:  1,
			},
			CrossDomain:    false,
			CircuitBreaker: &CircuitBreaker{Expression: "NetworkErrorRatio() > 0.5"},
			Recovery:       true,
		},
		Proxy: &Proxy{
			ProxyLog: &LogInfo{
				Path:  "./log/proxy.log",
				Level: "INFO",
			},
			Tls: &tlsx.Config{
				Ca:   "./build/package/proxy/cert/ca.pem",
				Cert: "./build/package/proxy/cert/client.pem",
				Key:  "./build/package/proxy/cert/client-key.pem",
			},
		},
	}
	file, err := os.Create("../../conf/config.json")
	require.NoError(t, err)
	defer file.Close()
	err = jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(file).Encode(cf)
	require.NoError(t, err)
}
