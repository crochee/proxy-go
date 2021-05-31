package config

import "github.com/crochee/proxy-go/config/dynamic"

type Proxy struct {
	ProxyLog *dynamic.LogInfo `json:"request_log,omitempty" yaml:"request_log,omitempty"`
	Tls      *TlsConfig       `json:"tls,omitempty" yaml:"tls,omitempty"`
}
