// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"time"

	"proxy-go/ptls"
)

type Config struct {
	Server *Server `yaml:"server,omitempty"`
}

type Server struct {
	Medata []*Medata `yaml:"medata,omitempty"`
}

type Medata struct {
	Name         string        `yaml:"name,omitempty"`
	Scheme       string        `yaml:"scheme,omitempty"`
	Port         int           `yaml:"port,omitempty"`
	Tls          *TlsConfig    `yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `yaml:"grace_time_out,omitempty"`
}

type TlsConfig struct {
	Cert ptls.FileOrContent `yaml:"cert,omitempty"`
	Key  ptls.FileOrContent `yaml:"key,omitempty"`
}
