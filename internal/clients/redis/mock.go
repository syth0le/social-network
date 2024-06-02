package redis

import (
	"context"
	"encoding"
	"fmt"

	xerrors "github.com/syth0le/gopnik/errors"
	"go.uber.org/zap"
)

type ClientMock struct {
	Logger *zap.Logger
}

func (c *ClientMock) HSet(ctx context.Context, hasTTL bool, key string, values ...any) error {
	c.Logger.Debug("hset through cache mock")
	return nil
}

func (c *ClientMock) HGet(ctx context.Context, key string, field string, scanTo encoding.BinaryUnmarshaler) error {
	c.Logger.Debug("hget through cache mock")
	return xerrors.WrapNotFoundError(fmt.Errorf("cannot find smth in cache mock"), "not found in cache")

}

func (c *ClientMock) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	c.Logger.Debug("hget through cache mock")
	return nil, xerrors.WrapNotFoundError(fmt.Errorf("cannot find smth in cache mock"), "not found in cache")
}

func (c *ClientMock) HSetNX(ctx context.Context, hasTTL bool, key string, field string, value encoding.BinaryMarshaler) error {
	c.Logger.Debug("hsetnx through cache mock")
	return nil
}

func (c *ClientMock) LPush(ctx context.Context, key string, value encoding.BinaryMarshaler) error {
	c.Logger.Debug("lpush through cache mock")
	return nil
}

func (c *ClientMock) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	c.Logger.Debug("lrange through cache mock")
	return nil, xerrors.WrapNotFoundError(fmt.Errorf("cannot find smth in cache mock"), "not found in cache")
}

func (c *ClientMock) Close() error {
	return nil
}
