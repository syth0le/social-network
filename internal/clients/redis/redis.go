package redis

import (
	"context"
	"encoding"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	xerrors "github.com/syth0le/gopnik/errors"
	"go.uber.org/zap"

	"social-network/cmd/social-network/configuration"
)

const (
	defaultClientName = "social-network"
)

type Client interface {
	HSet(ctx context.Context, hasTTL bool, key string, values ...any) error
	HGet(ctx context.Context, key string, field string, scanTo encoding.BinaryUnmarshaler) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSetNX(ctx context.Context, hasTTL bool, key string, field string, value encoding.BinaryMarshaler) error
	Close() error
}

type ClientImpl struct {
	Client             *redis.Client
	ExpirationDuration time.Duration
}

func NewRedisClient(logger *zap.Logger, cfg configuration.RedisConfig) Client {
	if !cfg.Enable {
		return &ClientMock{Logger: logger}
	}

	return &ClientImpl{
		Client: redis.NewClient(&redis.Options{
			Addr:       cfg.Address,
			ClientName: defaultClientName,
			Password:   cfg.Password,
			DB:         cfg.Database,
		}),
		ExpirationDuration: cfg.ExpirationDuration,
	}
}

func (c *ClientImpl) Close() error {
	return c.Client.Close()
}

func (c *ClientImpl) HSet(ctx context.Context, hasTTL bool, key string, values ...any) error {
	err := c.Client.HSet(ctx, key, values).Err()
	if err != nil {
		return err
	}

	if hasTTL {
		c.Client.Expire(ctx, key, c.ExpirationDuration)
	}

	return nil
}

func (c *ClientImpl) HGet(ctx context.Context, key string, field string, scanTo encoding.BinaryUnmarshaler) error {
	resp, err := c.Client.HGet(ctx, key, field).Result()
	if err != nil {
		if err != redis.Nil {
			return xerrors.WrapInternalError(fmt.Errorf("hget error"))
		}

		return xerrors.WrapNotFoundError(err, "not found in cache")
	}

	err = scanTo.UnmarshalBinary([]byte(resp))
	if err != nil {
		return fmt.Errorf("unmarshal error: %w", err)
	}

	return nil
}

func (c *ClientImpl) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	resp, err := c.Client.HGetAll(ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, xerrors.WrapInternalError(fmt.Errorf("hget error"))
		}

		return nil, xerrors.WrapNotFoundError(err, "not found in cache")
	}

	return resp, nil
}

func (c *ClientImpl) HSetNX(ctx context.Context, hasTTL bool, key string, field string, value encoding.BinaryMarshaler) error {
	err := c.Client.HSetNX(ctx, key, field, value).Err()
	if err != nil {
		return err
	}

	if hasTTL {
		c.Client.Expire(ctx, key, c.ExpirationDuration)
	}

	return nil
}
