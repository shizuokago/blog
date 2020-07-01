package config

import (
	"log"
)

type Config struct {
	Port        string
	DevelopMode bool
	ProjectID   string
}

var gConf Config

func Get() *Config {
	return &gConf
}

func Set(opts ...Option) error {
	gConf = Config{}
	for _, opt := range opts {
		err := opt(&gConf)
		if err != nil {
			log.Printf("%+v", err)
		}
	}
	return nil
}
