package redis

import (
	"context"
	"encoding"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	xerrors "github.com/syth0le/gopnik/errors"
	"go.uber.org/zap"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
)

const (
	defaultClientName = "social-network"
)

type Client interface {
	HSet(ctx context.Context, hasTTL bool, key string, values ...any) error
	HGet(ctx context.Context, key string, field string, scanTo encoding.BinaryUnmarshaler) error
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSetNX(ctx context.Context, hasTTL bool, key string, field string, value encoding.BinaryMarshaler) error
	LPush(ctx context.Context, key string, value encoding.BinaryMarshaler) error
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LRem(ctx context.Context, key string, value encoding.BinaryMarshaler) error
	Delete(ctx context.Context, keys ...string) error
	Close() error
}

type ClientImpl struct {
	Logger             *zap.Logger
	Client             *redis.Client
	ExpirationDuration time.Duration
	MaxListRange       int64
}

func NewRedisClient(logger *zap.Logger, cfg configuration.RedisConfig) Client {
	if !cfg.Enable {
		return &ClientMock{Logger: logger}
	}

	return &ClientImpl{
		Logger: logger,
		Client: redis.NewClient(&redis.Options{
			Addr:       cfg.Address,
			ClientName: defaultClientName,
			Password:   cfg.Password,
			DB:         cfg.Database,
		}),
		ExpirationDuration: cfg.ExpirationDuration,
		MaxListRange:       cfg.MaxListRange,
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

func (c *ClientImpl) LPush(ctx context.Context, key string, value encoding.BinaryMarshaler) error {
	_, err := c.Client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		_, err := pipe.LPush(ctx, key, value).Result()
		if err != nil {
			return fmt.Errorf("lpush: %w", err)
		}

		_, err = pipe.LTrim(ctx, key, 0, c.MaxListRange).Result()
		if err != nil {
			return fmt.Errorf("ltrim: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("tx pipelined: %w", err)
	}

	return nil
}

func (c *ClientImpl) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if start > c.MaxListRange {
		return nil, fmt.Errorf("start is greater than max list range value")
	}

	list, err := c.Client.LRange(ctx, key, max(start, 0), min(stop, c.MaxListRange-1)).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, xerrors.WrapInternalError(fmt.Errorf("lrange error"))
		}

		return nil, xerrors.WrapNotFoundError(err, "not found in list cache")
	}

	return list, nil
}

func (c *ClientImpl) LRem(ctx context.Context, key string, value encoding.BinaryMarshaler) error {
	_, err := c.Client.LRem(ctx, key, 0, value).Result()
	if err != nil {
		return fmt.Errorf("ltrim: %w", err)
	}

	return nil
}

func (c *ClientImpl) Delete(ctx context.Context, keys ...string) error {
	_, err := c.Client.Del(ctx, keys...).Result()
	if err != nil {
		return fmt.Errorf("del keys: %w", err)
	}

	return nil
}
