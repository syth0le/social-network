package configuration

import (
	"time"

	xstorage "github.com/syth0le/gopnik/db/postgres"
	xlogger "github.com/syth0le/gopnik/logger"
	xservers "github.com/syth0le/gopnik/servers"
)

const (
	defaultAppName   = "social-network"
	defaultSaltValue = "saltValue"
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
			GracefulShutdownTimeout: 15 * time.Second,
			ForceShutdownTimeout:    20 * time.Second,
			App:                     defaultAppName,
			SaltValue:               defaultSaltValue,
		},
		PublicServer: xservers.ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		AdminServer: xservers.ServerConfig{
			Enable:   false,
			Endpoint: "",
			Port:     0,
		},
		InternalGRPCServer: xservers.GRPCServerConfig{
			Port:             0,
			EnableRecover:    false,
			EnableReflection: false,
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
		Cache: RedisConfig{
			Enable:             false,
			Address:            "",
			Password:           "",
			Database:           0,
			ExpirationDuration: 5 * time.Minute,
			HeaterDuration:     24 * time.Hour,
			MaxListRange:       1000,
		},
		Queue: RabbitConfig{
			Enable:       false,
			Address:      "",
			QueueName:    "",
			ExchangeName: "",
		},
		NotificationsQueue: RabbitConfig{
			Enable:       false,
			Address:      "",
			QueueName:    "",
			ExchangeName: "",
		},
	}
}
