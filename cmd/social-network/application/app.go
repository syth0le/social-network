package application

import (
	"context"
	"fmt"
	"social-network/cmd/social-network/configuration"
	"social-network/internal/service/user"

	"go.uber.org/zap"
)

type App struct {
	Config *configuration.Config
	Logger *zap.Logger
	// Closer *closer.Closer // TODO -fix
}

func New(cfg *configuration.Config, logger *zap.Logger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
		// Closer: nil, // TODO nil
	}
}

func (a *App) Run() error {
	ctx, cancelFunction := context.WithCancel(context.Background())
	_ = cancelFunction

	envStruct, err := a.constructEnv(ctx)
	if err != nil {
		return fmt.Errorf("construct env: %w", err)
	}

	httpServer := a.newHTTPServer(envStruct)
	httpServer.Run()

	return nil
}

type env struct {
	userService user.Service
	// tokenService
}

func (a *App) constructEnv(ctx context.Context) (*env, error) {
	return nil, nil
}
