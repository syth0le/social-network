package application

import (
	"context"
	"fmt"
	"syscall"

	"go.uber.org/zap"

	xcloser "github.com/syth0le/gopnik/closer"

	"social-network/cmd/users-generator/configuration"
	"social-network/internal/service/generator"
	"social-network/internal/storage/postgres"
)

type App struct {
	Config *configuration.Config
	Logger *zap.Logger
	Closer *xcloser.Closer
}

func New(cfg *configuration.Config, logger *zap.Logger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
		Closer: xcloser.NewCloser(logger, cfg.Application.GracefulShutdownTimeout, cfg.Application.ForceShutdownTimeout, syscall.SIGINT, syscall.SIGTERM),
	}
}

func (a *App) Run() error {
	ctx, cancelFunction := context.WithCancel(context.Background())
	a.Closer.Add(func() error {
		cancelFunction()
		return nil
	})

	db, err := postgres.NewStorage(a.Logger, a.Config.Storage)
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}
	a.Closer.Add(db.Close)

	userService := &generator.ServiceImpl{
		Storage: db,
	}

	if err := userService.BatchGenerateUsers(ctx); err != nil {
		return fmt.Errorf("batch generate users: %w", err)
	}

	a.Closer.Wait()
	return nil
}
