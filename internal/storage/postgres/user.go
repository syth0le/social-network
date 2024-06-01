package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	xerrors "github.com/syth0le/gopnik/errors"

	"social-network/internal/model"
)

// TODO: tests

func (s *Storage) GetUserByLogin(ctx context.Context, userLogin *model.UserLogin) (*model.User, error) {
	sql, args, err := sq.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{
			fieldUsername:  userLogin.Username,
			fieldDeletedAt: nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
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
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) GetUserByID(ctx context.Context, id model.UserID) (*model.User, error) {
	// s.storage.logger.Sugar().Infof("some info for debug: %v", id) TODO make storage logger
	sql, args, err := sq.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{
			fieldID:        id.String(),
			fieldDeletedAt: nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) GetUserByTokenAndID(ctx context.Context, userID model.UserID, token string) (*model.User, error) {
	sql, args, err := sq.Select(tableFields(UserTable, userFields)...).
		From(UserTable).
		Join(joinString(UserTable, fieldID, TokenTable, fieldUserID)).
		Where(sq.Eq{
			tableField(UserTable, fieldID):         userID,
			tableField(TokenTable, fieldToken):     token,
			tableField(TokenTable, fieldDeletedAt): nil,
			tableField(UserTable, fieldDeletedAt):  nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity userEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) SearchUser(ctx context.Context, firstName, secondName string) ([]*model.User, error) {
	sql, args, err := sq.Select(userFields...).
		From(UserTable).
		Where(sq.And{
			sq.ILike{fieldFirstName: firstName + "%"},
			sq.ILike{fieldSecondName: secondName + "%"},
			sq.Eq{
				fieldDeletedAt: nil,
			},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entities []userEntity
	err = sqlx.SelectContext(ctx, s.Slave(), &entities, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return userEntitiesToModels(entities), nil
}

func (s *Storage) BatchCreateUser(ctx context.Context, users []*model.UserRegister) error {
	now := time.Now().Truncate(time.Millisecond)
	query := sq.Insert(UserTable).
		Columns(userFields...)

	for _, user := range users {
		query = query.Values(user.ID, user.Username, user.HashedPassword, user.FirstName, user.SecondName,
			user.Sex, user.Birthdate, user.Biography, user.City, now)
	}

	sql, args, err := query.
		Suffix(returningUser).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	_, err = s.Master().ExecContext(ctx, sql, args...)
	if err != nil {
		return xerrors.WrapSqlError(err)
	}

	return nil
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
		UserID:         model.UserID(entity.ID),
		Username:       entity.Username,
		HashedPassword: entity.HashedPassword,
		FirstName:      entity.FirstName,
		SecondName:     entity.SecondName,
		Sex:            entity.Sex,
		Birthdate:      entity.Birthdate,
		Biography:      entity.Biography,
		City:           entity.City,
	}
}

func userEntitiesToModels(entities []userEntity) []*model.User {
	var userModels []*model.User
	for _, entity := range entities {
		userModels = append(userModels, userEntityToModel(entity))
	}
	return userModels
}
