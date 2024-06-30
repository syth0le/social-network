package tarantool

import (
	"context"
	"fmt"
	"strconv"
	"time"

	xerrors "github.com/syth0le/gopnik/errors"
	"github.com/tarantool/go-tarantool/v2"
	"go.uber.org/zap"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
	"github.com/syth0le/social-network/internal/model"
)

type Storage struct {
	conn   *tarantool.Connection
	logger *zap.Logger
}

func NewStorage(ctx context.Context, logger *zap.Logger, config configuration.TarantoolConfig) (*Storage, error) {
	dialer := tarantool.NetDialer{
		Address:  config.Address,
		User:     config.Username,
		Password: config.Password,
	}
	opts := tarantool.Opts{
		Timeout:   config.TimeoutDuration,
		Reconnect: 5 * time.Second,
		RateLimit: 1000,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, fmt.Errorf("connection refused: %w", err)
	}

	return &Storage{
		conn:   conn,
		logger: logger,
	}, nil
}

func (s *Storage) AddUser(ctx context.Context, user *model.TarantoolUser) (tarantool.Response, error) {
	resp, err := s.conn.Do(
		tarantool.NewInsertRequest("users").Tuple(makeUserModelToTuple(user)).Context(ctx),
	).GetResponse()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("insert user: %w", err))
	}

	return resp, nil
}

func (s *Storage) SearchUser(ctx context.Context, firstName, lastName string) ([]model.TarantoolUser, error) {
	req := tarantool.NewCallRequest("search_by_first_second_name_with_size").Args([]interface{}{firstName, lastName, 10}).Context(ctx)
	resp, err := s.conn.Do(req).Get()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("search by firstname and secondname with offset: %w", err))
	}

	var userList []model.TarantoolUser
	for _, r := range resp {
		list := r.([]interface{})
		for _, rr := range list {
			rawUser := rr.([]interface{})
			userList = append(userList, makeUserModel(rawUser))
		}
	}

	return userList, err
}

func (s *Storage) Close() error {
	return s.conn.Close()
}

func makeUserModel(rawData []interface{}) model.TarantoolUser {
	return model.TarantoolUser{
		UserID:         "",
		FirstName:      rawData[1].(string),
		SecondName:     rawData[2].(string),
		Username:       rawData[3].(string),
		HashedPassword: rawData[4].(string),
		Sex:            rawData[5].(string),
		Biography:      rawData[6].(string),
		City:           rawData[7].(string),
	}
}

func makeUserModelToTuple(user *model.TarantoolUser) []interface{} {
	atoi, err := strconv.Atoi(user.UserID)
	if err != nil {
		return nil
	}
	return []any{atoi, user.FirstName, user.SecondName, user.Username, user.HashedPassword, user.Sex, user.Biography, user.City}
}
