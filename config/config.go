// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Cfg *Config

// InitConfig init Config
func InitConfig(path string) {
	configYaml, err := loadYaml(path)
	if err != nil {
		panic(err)
	}
	Cfg = configYaml
}

func loadYaml(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
