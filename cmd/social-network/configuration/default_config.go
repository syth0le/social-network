package configuration

import "social-network/internal/utils"

const (
	defaultAppName     = "social-network"
	defaultEnvironment = "dev"
)

func NewDefaultConfig() *Config {
	return &Config{
		Logger: &LoggerConfig{
			Level:       utils.InfoLevel,
			Encoding:    "console",
			Path:        "stdout",
			App:         defaultAppName,
			Environment: defaultEnvironment,
		},
		Server: &ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		Storage: &StorageConfig{
			EnableMock:         false,
			Host:               "",
			Port:               0,
			Database:           "",
			Username:           "",
			Password:           "",
			SSLMode:            "",
			ConnectionAttempts: 0,
		},
	}
}
