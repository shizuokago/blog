package config

import (
	"cloud.google.com/go/compute/metadata"
	"golang.org/x/xerrors"
)

type Config struct {
	Port        string
	DevelopMode bool
	ProjectID   string
}

var gConf *Config

func init() {
	gConf = defaultConfig()
}

func defaultConfig() *Config {
	var conf Config
	conf.DevelopMode = !metadata.OnGCE()
	conf.Port = "8080"
	conf.ProjectID = "blog"
	return &conf
}

func Get() *Config {
	return gConf
}

func Set(opts ...Option) error {
	for idx, opt := range opts {
		err := opt(gConf)
		if err != nil {
			return xerrors.Errorf("Option[%d] error: %w", idx, err)
		}
	}
	return nil
}
