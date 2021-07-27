package dynamic

import "github.com/crochee/proxy/pkg/tlsx"

type Proxy struct {
	ProxyLog *LogInfo     `json:"request_log,omitempty" yaml:"request_log,omitempty"`
	Tls      *tlsx.Config `json:"tls,omitempty" yaml:"tls,omitempty"`
}
