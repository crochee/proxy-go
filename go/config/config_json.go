package config

import (
	"os"

	"github.com/json-iterator/go"

	"github.com/crochee/proxy/internal"
)

type Json struct {
	path string
}

func (j Json) Decode() (*Spec, error) {
	file, err := os.Open(j.path)
	if err != nil {
		return nil, err
	}
	defer internal.Close(file)
	var config Spec
	dec := jsoniter.ConfigCompatibleWithStandardLibrary.NewDecoder(file)
	dec.UseNumber() // 解决json 将int当成float的情况
	if err = dec.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (j Json) Encode(c *Spec) error {
	file, err := os.Create(j.path)
	if err != nil {
		return err
	}
	defer internal.Close(file)
	return jsoniter.ConfigCompatibleWithStandardLibrary.NewEncoder(file).Encode(c)
}
