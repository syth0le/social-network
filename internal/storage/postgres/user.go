package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"social-network/internal/model"
	"social-network/internal/utils"
	"time"
)

// TODO: tests

func (s *Storage) LoginUser(ctx context.Context, userLogin *model.UserLogin) (*model.User, error) {
	sql, args, err := sq.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{
			fieldUsername:       userLogin.Username,
			fieldHashedPassword: userLogin.HashedPassword,
			fieldDeletedAt:      nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) CreateUser(ctx context.Context, params *model.UserRegister) (*model.User, error) {
	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("params validate: %w", err)
	}
	now := time.Now().Truncate(time.Millisecond)
	sql, args, err := sq.Insert(UserTable).
		Columns(userFields...).
		Values(params.ID, params.Username, params.HashedPassword, params.FirstName, params.SecondName,
			params.Sex, params.Birthdate, params.Biography, params.City, now,
		).
		Suffix(returningUser).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}
	var entity userEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) GetUserByID(ctx context.Context, id model.UserID) (*model.User, error) {
	s.storage.logger.Sugar().Infof("some info for debug: %v", id)
	sql, args, err := sq.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{
			fieldID:        id.String(),
			fieldDeletedAt: nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

type userEntity struct {
	ID             string    `db:"id"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hashed_password"`
	FirstName      string    `db:"first_name"`
	SecondName     string    `db:"second_name"`
	Sex            string    `db:"sex"`
	Birthdate      time.Time `db:"birthdate"`
	Biography      string    `db:"biography"`
	City           string    `db:"city"`
	CreatedAt      time.Time `db:"created_at"`
}

func userEntityToModel(entity userEntity) *model.User {
	return &model.User{
		UserID:     model.UserID(entity.ID),
		Username:   entity.Username,
		FirstName:  entity.FirstName,
		SecondName: entity.SecondName,
		Sex:        entity.Sex,
		Birthdate:  entity.Birthdate,
		Biography:  entity.Biography,
		City:       entity.City,
	}
}
