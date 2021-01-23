// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/23

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Yml struct {
	path string
}

func (y Yml) Decode() (*Config, error) {
	file, err := os.Open(y.path)
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

func (y Yml) Encode(c *Config) error {
	file, err := os.Create(y.path)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewEncoder(file).Encode(c)
}
