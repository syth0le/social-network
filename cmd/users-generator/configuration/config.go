package configuration

import (
	xstorage "github.com/syth0le/gopnik/db/postgres"
	xlogger "github.com/syth0le/gopnik/logger"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
)

type Config struct {
	Logger      xlogger.LoggerConfig          `yaml:"logger"`
	Application ApplicationConfig             `yaml:"application"`
	Storage     xstorage.StorageConfig        `yaml:"storage"`
	Tarantool   configuration.TarantoolConfig `yaml:"tarantool"`
}

func (c *Config) Validate() error {
	return nil // todo
}

type ApplicationConfig struct {
	App      string `yaml:"app"`
	DataFile string `yaml:"data_file"`
}

func (c *ApplicationConfig) Validate() error {
	return nil // todo
}
