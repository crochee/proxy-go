// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package dynamic

type Middleware struct {
	Prefix *prefix `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

type prefix struct {
	Path []string `json:"path,omitempty" yaml:"path,omitempty"`
}
