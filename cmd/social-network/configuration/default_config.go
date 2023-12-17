package configuration

import (
	"social-network/internal/utils"
	"time"
)

const (
	defaultAppName   = "social-network"
	defaultSaltValue = "saltValue"
)

func NewDefaultConfig() *Config {
	return &Config{
		Logger: LoggerConfig{
			Level:       utils.InfoLevel,
			Encoding:    "console",
			Path:        "stdout",
			Environment: utils.Development,
		},
		Application: ApplicationConfig{
			GracefulShutdownTimeout: 15 * time.Second,
			App:                     defaultAppName,
			SaltValue:               defaultSaltValue,
		},
		PublicServer: ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		AdminServer: ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		Storage: StorageConfig{
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
	}
}
