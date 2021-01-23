// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/23

package dynamic

import (
	"time"

	"proxy-go/ptls"
)

type Medata struct {
	Name         string        `json:"name,omitempty" yaml:"name,omitempty"`
	Scheme       string        `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Port         int           `json:"port,omitempty" yaml:"port,omitempty"`
	Tls          *TlsConfig    `json:"tls,omitempty" yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `json:"grace_time_out,omitempty" yaml:"grace_time_out,omitempty"`
	Mode         string        `json:"mode,omitempty" yaml:"mode,omitempty"`
	Path         string        `json:"path,omitempty" yaml:"path,omitempty"`
	LocationList []*Location   `json:"location_list,omitempty" yaml:"location_list,omitempty"`
	Middleware   *Middleware   `json:"middleware,omitempty" yaml:"middleware,omitempty"`
}

type TlsConfig struct {
	Cert ptls.FileOrContent `json:"cert,omitempty" yaml:"cert,omitempty"`
	Key  ptls.FileOrContent `json:"key,omitempty" yaml:"key,omitempty"`
}

type Location struct {
	ProxyPass string `json:"proxy_pass,omitempty" yaml:"proxy_pass,omitempty"`
	Weight    int    `json:"weight,omitempty" yaml:"weight,omitempty"`
}
