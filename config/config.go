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
	"github.com/crochee/proxy-go/ptls"
)

type Config struct {
	*Spec
	lc LoadConfig
}
type Spec struct {
	Medata     *Medata         `json:"medata" yaml:"medata"`
	Middleware *dynamic.Config `json:"middleware,omitempty" yaml:"middleware,omitempty"`
}

type Medata struct {
	Tls          *TlsConfig    `json:"tls,omitempty" yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `json:"grace_time_out,omitempty" yaml:"grace_time_out,omitempty"`

	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host   string `json:"host" yaml:"host"`

	LogPath  string `json:"log_path,omitempty" yaml:"log_path,omitempty"`
	LogLevel string `json:"log_level,omitempty" yaml:"log_level,omitempty"`
}

type TlsConfig struct {
	Ca   ptls.FileOrContent `json:"cert" yaml:"cert"`
	Cert ptls.FileOrContent `json:"cert" yaml:"cert"`
	Key  ptls.FileOrContent `json:"key" yaml:"key"`
}

var Cfg *Config

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

func loadConfig(path string) (*Config, error) {
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
	return &Config{
		Spec: spec,
		lc:   lc,
	}, nil
}
