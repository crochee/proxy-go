// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/crochee/proxy-go/config/dynamic"
	"github.com/crochee/proxy-go/internal/tlsx"
)

type Spec struct {
	Medata     *Medata         `json:"medata" yaml:"medata"`
	Middleware *dynamic.Config `json:"middleware,omitempty" yaml:"middleware,omitempty"`
	Proxy      *TlsConfig      `json:"proxy,omitempty" yaml:"proxy,omitempty"`
}

type Medata struct {
	Tls          *TlsConfig    `json:"tls,omitempty" yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `json:"grace_time_out,omitempty" yaml:"grace_time_out,omitempty"`

	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host   string `json:"host" yaml:"host"`

	SystemLog  *dynamic.LogInfo `json:"system_log,omitempty" yaml:"system_log,omitempty"`
	RequestLog *dynamic.LogInfo `json:"request_log,omitempty" yaml:"request_log,omitempty"`
}

type TlsConfig struct {
	Ca   tlsx.FileOrContent `json:"ca" yaml:"ca"`
	Cert tlsx.FileOrContent `json:"cert" yaml:"cert"`
	Key  tlsx.FileOrContent `json:"key" yaml:"key"`
}

var Cfg *Spec

// InitConfig init Config
func InitConfig(path string) {
	config, err := loadConfig(path)
	if err != nil {
		panic(err)
	}
	Cfg = config
}

type LoadConfig interface {
	Decode() (*Spec, error)
	Encode(*Spec) error
}

func loadConfig(path string) (*Spec, error) {
	var lc LoadConfig
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".json":
		lc = Json{path: path}
	case ".yml", ".yaml":
		lc = Yml{path: path}
	default:
		return nil, fmt.Errorf("unsupport config extension %s", ext)
	}
	spec, err := lc.Decode()
	if err != nil {
		return nil, err
	}
	return spec, nil
}
