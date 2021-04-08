package config

import (
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	"golang.org/x/xerrors"
)

type Option func(*Config) error

func Port() Option {
	return func(c *Config) error {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		c.Port = port
		return nil
	}
}

const (
	DatastoreEmulatorHostEnv     = "DATASTORE_EMULATOR_HOST"
	DatastoreProjectIDEnv        = "DATASTORE_PROJECT_ID"
	DatastoreDatasetEnv          = "DATASTORE_DATASET"
	DefaultDatastoreEmulatorHost = "localhost:8081"
)

func Project() Option {

	return func(c *Config) error {

		c.DevelopMode = true
		c.ProjectID = "blog"

		//2020/7/1 現在AppEngine実行パスで判定
		wd, err := os.Getwd()
		if err != nil {
			return xerrors.Errorf("get work directory : %w", err)
		}

		if wd == "/srv" {
			c.DevelopMode = false
		}

		if !c.DevelopMode {
			c.ProjectID, err = metadata.ProjectID()
			if err != nil {
				return xerrors.Errorf("get project id: %w", err)
			}
		}

		fmt.Println("ProjectID=" + c.ProjectID)
		return nil
	}

}

func Datastore() Option {

	return func(c *Config) error {

		if c.DevelopMode {
			host := os.Getenv(DatastoreEmulatorHostEnv)
			if host == "" {
				host = DefaultDatastoreEmulatorHost
				os.Setenv(DatastoreEmulatorHostEnv, DefaultDatastoreEmulatorHost)
			}

			fmt.Println("Develop DatastoreHost=" + host)

			if os.Getenv(DatastoreProjectIDEnv) == "" {
				os.Setenv(DatastoreProjectIDEnv, c.ProjectID)
			}
		}

		if os.Getenv(DatastoreDatasetEnv) == "" {
			os.Setenv(DatastoreDatasetEnv, c.ProjectID)
		}

		return nil
	}
}

func ProjectID() string {
	return gConf.ProjectID
}
