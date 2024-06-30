package configuration

import (
	"time"

	xstorage "github.com/syth0le/gopnik/db/postgres"
	xlogger "github.com/syth0le/gopnik/logger"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
)

const (
	defaultAppName = "users-generator"
)

func NewDefaultConfig() *Config {
	return &Config{
		Logger: xlogger.LoggerConfig{
			Level:       xlogger.InfoLevel,
			Encoding:    "console",
			Path:        "stdout",
			Environment: xlogger.Development,
		},
		Application: ApplicationConfig{
			App:      defaultAppName,
			DataFile: "",
		},
		Storage: xstorage.StorageConfig{
			EnableMock:            false,
			Hosts:                 []string{},
			Port:                  0,
			Database:              "",
			Username:              "",
			Password:              "",
			SSLMode:               "",
			ConnectionAttempts:    0,
			InitializationTimeout: 5 * time.Second,
		},
		Tarantool: configuration.TarantoolConfig{
			EnableMock:      false,
			Address:         "",
			Username:        "",
			Password:        "",
			TimeoutDuration: 1 * time.Second,
		},
	}
}
