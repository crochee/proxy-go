// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/23

package config

import (
	"os"

	"github.com/json-iterator/go"
)

type Json struct {
	path string
}

func (j Json) Decode() (*Config, error) {
	file, err := os.Open(j.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	if err = jsoniter.ConfigFastest.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (j Json) Encode(c *Config) error {
	file, err := os.Create(j.path)
	if err != nil {
		return err
	}
	defer file.Close()
	return jsoniter.ConfigFastest.NewEncoder(file).Encode(c)
}
