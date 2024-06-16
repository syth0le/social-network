package application

import (
	"context"
	"fmt"
	"syscall"

	"go.uber.org/zap"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
	"github.com/syth0le/social-network/internal/authentication"
	"github.com/syth0le/social-network/internal/clients/rabbit"
	"github.com/syth0le/social-network/internal/clients/redis"
	"github.com/syth0le/social-network/internal/infrastructure_services/cache"
	"github.com/syth0le/social-network/internal/infrastructure_services/heater"
	"github.com/syth0le/social-network/internal/infrastructure_services/queue"
	"github.com/syth0le/social-network/internal/service/friend"
	"github.com/syth0le/social-network/internal/service/post"
	"github.com/syth0le/social-network/internal/service/user"
	"github.com/syth0le/social-network/internal/storage/postgres"
	"github.com/syth0le/social-network/internal/token"

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

	internalGRPCServer := a.newInternalGRPCServer(envStruct)
	a.Closer.AddForce(internalGRPCServer.ForcefullyStop)
	a.Closer.Add(internalGRPCServer.GracefullyStop)

	a.Closer.Run(internalGRPCServer.Run)

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

	redisClient := redis.NewRedisClient(a.Logger, a.Config.Cache) // TODO: move to gopnik
	a.Closer.Add(redisClient.Close)

	cacheService := &cache.ServiceImpl{Client: redisClient, Logger: a.Logger}

	publisher, err := rabbit.NewRabbitPublisher(a.Logger, a.Config.Queue) // TODO: move to gopnik
	if err != nil {
		return nil, fmt.Errorf("new rabbit publisher: %w", err)
	}
	a.Closer.Add(publisher.Close)

	consumer, err := rabbit.NewRabbitConsumer(a.Logger, a.Config.Queue) // TODO: move to gopnik
	if err != nil {
		return nil, fmt.Errorf("new rabbit consumer: %w", err)
	}
	a.Closer.Add(consumer.Close)

	tokenManager := token.NewManager(a.Config.Application)
	userService := &user.ServiceImpl{
		Storage:      postgresDB,
		TokenManager: tokenManager,
	}

	friendService := &friend.ServiceImpl{Storage: postgresDB}

	producerService := &queue.ProducerService{
		Logger:   a.Logger,
		Producer: publisher,
	}

	postService := &post.ServiceImpl{
		Storage:         postgresDB,
		Cache:           cacheService,
		Logger:          a.Logger,
		ProducerService: producerService,
	}

	consumerService := queue.ConsumerService{
		Logger:        a.Logger,
		Consumer:      consumer,
		FriendService: friendService,
		CacheService:  cacheService,
	}
	a.Closer.Run(consumerService.Run)

	cacheHeater := heater.HeaterService{
		Logger:         a.Logger,
		CacheService:   cacheService,
		PostService:    postService,
		FriendService:  friendService,
		HeaterDuration: a.Config.Cache.HeaterDuration,
	}
	err = cacheHeater.Run(ctx) // todo: make http admin handler for heater
	if err != nil {
		return nil, fmt.Errorf("cannot heating cache: %w", err)
	}

	return &env{
		userService: userService,
		authenticationService: authentication.Service{
			UserService:  userService,
			TokenManager: tokenManager,
			Logger:       a.Logger,
		},
		postService:   postService,
		friendService: friendService,
	}, nil
}
