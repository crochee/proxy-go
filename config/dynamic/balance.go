// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type BalanceConfig struct {
	RegisterApis []*ServiceApi     `json:"register_apis,omitempty" yaml:"register_apis,omitempty"`
	Transfers    []*ServiceBalance `json:"transfers,omitempty" yaml:"transfers,omitempty"`
}

type Node struct {
	Scheme   string            `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host     string            `json:"host"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Weight   float64           `json:"weight"`
}

type Balance struct {
	Selector string  `json:"selector" yaml:"selector"`
	Nodes    []*Node `json:"nodes" yaml:"nodes"`
}

type ServiceApi struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	Path        string `json:"path" yaml:"path"`
	Method      string `json:"method" yaml:"method"`
}

type ServiceBalance struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	Balance
}
