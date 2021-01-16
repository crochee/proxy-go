// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

type Config struct {
	Server Server `yaml:"server,omitempty"`
}

type Server struct {
	Port *Port `yaml:"port,omitempty"`
}

type Port struct {
	Http  int `yaml:"http,omitempty"`
	Https int `yaml:"https,omitempty"`
}
