package application

import (
	"social-network/cmd/social-network/configuration"

	"go.uber.org/zap"
)

type App struct {
	Config *configuration.Config
	Logger *zap.Logger
	Closer *closer.Closer // TODO
}

func New(cfg *configuration.Config, logger *zap.Logger) *App {
	return &App{
		Config: cfg,
		Logger: logger,
		Closer: nil, // TODO nil
	}
}

func (a *App) Run() error {
	return nil
}
