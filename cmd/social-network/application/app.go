package application

import (
	"context"
	"fmt"
	"syscall"

	"go.uber.org/zap"

	"social-network/cmd/social-network/configuration"
	"social-network/internal/authentication"
	"social-network/internal/service/user"
	"social-network/internal/storage/postgres"
	"social-network/internal/token"

	xcloser "github.com/syth0le/gopnik/closer"
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

	envStruct, err := a.constructEnv(ctx)
	if err != nil {
		return fmt.Errorf("construct env: %w", err)
	}

	httpServer := a.newHTTPServer(envStruct)
	a.Closer.Add(httpServer.GracefulStop()...)

	a.Closer.Run(httpServer.Run()...)
	a.Closer.Wait()
	return nil
}

type env struct {
	userService           user.Service
	authenticationService authentication.Service
}

func (a *App) constructEnv(ctx context.Context) (*env, error) {
	db, err := postgres.NewStorage(a.Logger, a.Config.Storage)
	if err != nil {
		return nil, fmt.Errorf("new storage: %w", err)
	}
	a.Closer.Add(db.Close)

	tokenGenerator := token.NewGenerator(a.Config.Application)
	userService := &user.ServiceImpl{
		Storage:        db,
		TokenGenerator: tokenGenerator,
	}

	return &env{
		userService:           userService,
		authenticationService: authentication.Service{UserService: userService},
	}, nil
}
