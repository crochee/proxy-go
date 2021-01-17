// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/1/17

package config

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestInitConfig(t *testing.T) {
	cf := &Config{
		Server: &Server{
			Medata: []*Medata{
				{
					Name:   "server1",
					Scheme: "http",
					Port:   8079,
				},
				{
					Name:   "server2",
					Scheme: "https",
					Port:   8080,
				},
			},
		},
	}

	file, err := os.Create("../conf/config.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err = yaml.NewEncoder(file).Encode(cf); err != nil {
		t.Fatal(err)
	}
}
