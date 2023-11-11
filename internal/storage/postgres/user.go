package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"social-network/internal/model"
	"time"
)

func (s *Storage) LoginUser(ctx context.Context, userLogin *model.UserLogin) (*model.User, error) {

	return nil, nil
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
		return nil, fmt.Errorf("incorrect sql") // todo internal error
	}
	var entity userEntity
	err = sqlx.GetContext(ctx, nil, &entity, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	return userEntityToModel(entity), nil
}

func (s *Storage) GetUserByID(ctx context.Context, id model.UserID) (*model.User, error) {
	return nil, nil
}

type userEntity struct {
	ID         string    `db:"id"`
	Username   string    `db:"id"`
	FirstName  string    `db:"id"`
	SecondName string    `db:"id"`
	Age        int       `db:"id"`
	Sex        string    `db:"id"`
	Birthdate  time.Time `db:"id"`
	Biography  string    `db:"id"`
	City       string    `db:"id"`
}

func userEntityToModel(entity userEntity) *model.User {
	return &model.User{
		UserID:     model.UserID(entity.ID),
		Username:   entity.Username,
		FirstName:  entity.FirstName,
		SecondName: entity.SecondName,
		Age:        entity.Age,
		Sex:        entity.Sex,
		Birthdate:  entity.Birthdate,
		Biography:  entity.Biography,
		City:       entity.City,
	}
}
