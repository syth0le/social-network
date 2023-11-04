package configuration

import "social-network/internal/utils"

const (
	defaultAppName    = "social-network"
	defaultEnviroment = "dev"
)

func NewDefaultConfig() *Config {
	return &Config{
		Logger: &LoggerConfig{
			Level:      utils.InfoLevel,
			Encoding:   "console",
			Path:       "stdout",
			App:        defaultAppName,
			Enviroment: defaultEnviroment,
		},
		Server: &ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     "",
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
