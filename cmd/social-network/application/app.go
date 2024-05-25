package application

import (
	"context"
	"fmt"
	"syscall"

	"go.uber.org/zap"

	"social-network/cmd/social-network/configuration"
	"social-network/internal/authentication"
	"social-network/internal/service/friend"
	"social-network/internal/service/post"
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
	postService           post.Service
	friendService         friend.Service
}

func (a *App) constructEnv(ctx context.Context) (*env, error) {
	postgresDB, err := postgres.NewStorage(a.Logger, a.Config.Storage)
	if err != nil {
		return nil, fmt.Errorf("new storage: %w", err)
	}
	a.Closer.Add(postgresDB.Close)

	// redisDB, err := redis.NewStorage(a.Logger, a.Config.Storage) // TODO: move to gopnik
	// redisClient := clients.NewRedisClient(a.Logger, a.Config.Redis)
	// if err != nil {
	//	return nil, fmt.Errorf("new redis client: %w", err)
	// }
	// a.Closer.Add(redisClient.Close)

	// redisDB, err := kafka.NewProducer(a.Logger, a.Config.Storage) // TODO: move to gopnik
	// kafkaConsumer := clients.NewKafkaConsumer(a.Logger, a.Config.Redis)
	// if err != nil {
	//	return nil, fmt.Errorf("new redis client: %w", err)
	// }
	// a.Closer.Add(kafkaConsumer.Close)

	tokenManager := token.NewManager(a.Config.Application)
	userService := &user.ServiceImpl{
		Storage:      postgresDB,
		TokenManager: tokenManager,
	}

	return &env{
		userService: userService,
		authenticationService: authentication.Service{
			UserService:  userService,
			TokenManager: tokenManager,
			Logger:       a.Logger,
		},
		postService:   &post.ServiceImpl{Storage: postgresDB}, // TODO: add redis
		friendService: &friend.ServiceImpl{Storage: postgresDB},
	}, nil
}
