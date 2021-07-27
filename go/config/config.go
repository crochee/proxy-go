package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/pkg/tlsx"
)

type Spec struct {
	Medata *Medata `json:"medata" yaml:"medata"`
	Server *Medata `json:"server" yaml:"server"`
}

type Medata struct {
	Tls          *tlsx.Config  `json:"tls,omitempty" yaml:"tls,omitempty"`
	GraceTimeOut time.Duration `json:"grace_time_out,omitempty" yaml:"grace_time_out,omitempty"`

	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host   string `json:"host" yaml:"host"`

	SystemLog  *dynamic.LogInfo `json:"system_log,omitempty" yaml:"system_log,omitempty"`
	RequestLog *dynamic.LogInfo `json:"request_log,omitempty" yaml:"request_log,omitempty"`
}

var Cfg *Spec

// InitConfig init Config
func InitConfig(path string) error {
	config, err := loadConfig(path)
	if err != nil {
		return err
	}
	Cfg = config
	return nil
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
