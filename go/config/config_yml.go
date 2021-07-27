package config

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/crochee/proxy/internal"
)

type Yml struct {
	path string
}

func (y Yml) Decode() (*Spec, error) {
	file, err := os.Open(y.path)
	if err != nil {
		return nil, err
	}
	defer internal.Close(file)
	var config Spec
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (y Yml) Encode(c *Spec) error {
	file, err := os.Create(y.path)
	if err != nil {
		return err
	}
	defer internal.Close(file)
	return yaml.NewEncoder(file).Encode(c)
}
