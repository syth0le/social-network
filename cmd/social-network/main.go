package main

import (
	"context"
	"fmt"
	"log"
	"social-network/cmd/social-network/application"
	"social-network/cmd/social-network/configuration"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	cfg, err := loadConfig(ctx)
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("config validation failed: %v", err)
	}

	logger, err := constructLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	app := application.New(cfg, logger) // TODO: closures
	if err = app.Run(); err != nil {
		logger.Sugar().Fatalf("application stopped with error: %v", err)
	} else {
		logger.Info("application stopped")
	}

}

func constructLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction() // TODO
	if err != nil {
		return nil, err
	}

	defer logger.Sync()
	return logger, nil
}

func loadConfig(ctx context.Context) (*configuration.Config, error) {
	cfg := configuration.NewDefaultConfig()

	configPath := pflag.StringP("config", "c", "", "config path")

	if err := cleanenv.ReadConfig(*configPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}
	return cfg, nil
}
