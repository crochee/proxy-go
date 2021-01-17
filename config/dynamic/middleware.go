// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package dynamic

type Middleware struct {
	ReplaceHost *ReplaceHost `yaml:"replaceHost,omitempty"`
}

type ReplaceHost struct {
	Scheme string `yaml:"scheme,omitempty"`
	Host   string `yaml:"host,omitempty"`
}
