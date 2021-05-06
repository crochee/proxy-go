// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/2/6

package dynamic

type Config struct {
	RateLimit *RateLimit `json:"rate_limit,omitempty" yaml:"rate_limit,omitempty"`

	Balance map[string]*Balance `json:"balance,omitempty" yaml:"balance,omitempty"`
}
