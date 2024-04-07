package configuration

import (
	"time"

	xstorage "github.com/syth0le/gopnik/db/postgres"
	xlogger "github.com/syth0le/gopnik/logger"
)

type Config struct {
	Logger      xlogger.LoggerConfig   `yaml:"logger"`
	Application ApplicationConfig      `yaml:"application"`
	Storage     xstorage.StorageConfig `yaml:"storage"`
}

func (c *Config) Validate() error {
	return nil // todo
}

type ApplicationConfig struct {
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout"`
	ForceShutdownTimeout    time.Duration `yaml:"force_shutdown_timeout"`
	App                     string        `yaml:"app"`
}

func (c *ApplicationConfig) Validate() error {
	return nil // todo
}
