package config

import (
	"os"
)

type Config struct {
	Port string
}

var gConf Config

func Get() *Config {
	return &gConf
}

func Set(opts ...Option) error {
	gConf = Config{}

	for _, opt := range opts {
		opt(&gConf)
	}
	return nil
}

type Option func(*Config)

func AppEnginePort() Option {
	return func(c *Config) {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		c.Port = port
	}
}
