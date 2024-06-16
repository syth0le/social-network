package application

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/syth0le/social-network/cmd/users-generator/configuration"
	"github.com/syth0le/social-network/internal/service/generator"
	"github.com/syth0le/social-network/internal/storage/postgres"
)

type App struct {
	Config *configuration.Config
	Logger *zap.Logger
}

func New(cfg *configuration.Config, logger *zap.Logger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
	}
}

func (a *App) Run() error {
	ctx := context.Background()

	db, err := postgres.NewStorage(a.Logger, a.Config.Storage)
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}
	defer db.Close()

	userService := &generator.ServiceImpl{
		Logger:   a.Logger,
		Storage:  db,
		DataFile: a.Config.Application.DataFile,
	}

	if err := userService.BatchGenerateUsers(ctx); err != nil {
		return fmt.Errorf("batch generate users: %w", err)
	}

	return nil
}
