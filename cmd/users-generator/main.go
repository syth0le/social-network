package main

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
	xlogger "github.com/syth0le/gopnik/logger"
	"go.uber.org/zap"

	"social-network/cmd/users-generator/application"
	"social-network/cmd/users-generator/configuration"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	if err = cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	logger, err := constructLogger(cfg.Logger)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	app := application.New(cfg, logger)
	if err = app.Run(); err != nil {
		logger.Sugar().Fatalf("application stopped with error: %v", err)
	} else {
		logger.Info("application stopped")
	}
}

func constructLogger(cfg xlogger.LoggerConfig) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error
	switch cfg.Environment {
	case xlogger.Development:
		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("new development logger")
		}
	case xlogger.Production:
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("new production logger")
		}
	default:
		return nil, fmt.Errorf("unexpected environment for logger: %w", err)
	}

	defer logger.Sync()
	return logger, nil
}

func loadConfig() (*configuration.Config, error) {
	cfg := configuration.NewDefaultConfig()

	configPath := pflag.StringP("config", "c", "", "config path")
	pflag.Parse()

	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}
	return cfg, nil
}
