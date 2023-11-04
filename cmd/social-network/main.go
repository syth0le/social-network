package main

import (
	"context"
	"fmt"
	"log"
	"social-network/cmd/social-network/configuration"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/pflag"
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

	// TODO: closures

}

func loadConfig(ctx context.Context) (*configuration.Config, error) {
	cfg := configuration.NewDefaultConfig()

	configPath := pflag.StringP("config", "c", "", "config path")

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}
	return cfg, nil
}
