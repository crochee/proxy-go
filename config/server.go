// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/23

package config

import (
	"time"

	"proxy-go/ptls"
)

type Medata struct {
	Tls          *TlsConfig    `json:"tls,omitempty" yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `json:"grace_time_out,omitempty" yaml:"grace_time_out,omitempty"`

	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Port   int    `json:"port,omitempty" yaml:"port,omitempty"`
}

type TlsConfig struct {
	Cert ptls.FileOrContent `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  ptls.FileOrContent `json:"key,omitempty" yaml:"key,omitempty"`
}
