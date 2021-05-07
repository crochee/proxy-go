// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type Config struct {
	Balance     map[string]*Balance `json:"balance,omitempty" yaml:"balance,omitempty"`
	AccessLog   *LogInfo            `json:"access_log,omitempty" yaml:"access_log,omitempty"`
	RateLimit   *RateLimit          `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`
	Recovery    bool                `json:"recovery,omitempty" yaml:"recovery,omitempty"`
	CrossDomain bool                `json:"cross_domain,omitempty" yaml:"cross_domain,omitempty"`
}
