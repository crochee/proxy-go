// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type Config struct {
	BalanceNode map[string]*BalanceNode `json:"balance_node,omitempty" yaml:"balance_node,omitempty"`
	RateLimit   *RateLimit              `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
}
