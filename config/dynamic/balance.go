// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type BalanceNode struct {
	Scheme   string            `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host     string            `json:"host"`
	Metadata map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Weight   float64           `json:"weight"`
}
