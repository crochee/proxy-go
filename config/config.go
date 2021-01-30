// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"fmt"
	"strings"

	"path/filepath"
)

type Config struct {
	Server *Server `json:"server,omitempty" yaml:"server,omitempty"`
}

type Server struct {
	Medata []*Medata `json:"medata,omitempty" yaml:"medata,omitempty"`
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
	Decode() (*Config, error)
	Encode(*Config) error
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
	return lc.Decode()
}
