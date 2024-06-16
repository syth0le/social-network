package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	xerrors "github.com/syth0le/gopnik/errors"

	"github.com/syth0le/social-network/internal/model"
)

func (s *Storage) GetFriend(ctx context.Context, authorID, followerID model.UserID) (*model.Friend, error) {
	sql, args, err := sq.Select(
		tableField(UserTable, fieldID),
		tableField(UserTable, fieldUsername),
		tableField(UserTable, fieldFirstName),
		tableField(UserTable, fieldSecondName),
	).From(FriendTable).
		Join(
			joinString(FriendTable, fieldSecondUserID, UserTable, fieldID),
		).
		Where(sq.Eq{
			tableField(FriendTable, fieldFirstUserID):  authorID,
			tableField(FriendTable, fieldSecondUserID): followerID,
			tableField(FriendTable, fieldDeletedAt):    nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity friendEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return friendEntityToModel(entity), nil
}

func (s *Storage) ListFriends(ctx context.Context, userID model.UserID) ([]*model.Friend, error) {
	sql, args, err := sq.Select(
		tableField(UserTable, fieldID),
		tableField(UserTable, fieldUsername),
		tableField(UserTable, fieldFirstName),
		tableField(UserTable, fieldSecondName),
	).From(FriendTable).
		Join(
			joinString(FriendTable, fieldSecondUserID, UserTable, fieldID),
		).
		Where(sq.Eq{
			tableField(FriendTable, fieldFirstUserID): userID,
			tableField(FriendTable, fieldDeletedAt):   nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entities []friendEntity
	err = sqlx.SelectContext(ctx, s.Slave(), &entities, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return friendEntitiesToModels(entities), nil
}

func (s *Storage) AddFriend(ctx context.Context, params *model.AddFriendParams) error {
	err := params.Validate()
	if err != nil {
		return fmt.Errorf("params validate: %w", err)
	}

	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Insert(FriendTable).
		Columns(friendFields...).
		Values(params.FID, params.AuthorID, params.FollowerID, now).
		Values(params.SID, params.FollowerID, params.AuthorID, now).
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

func (s *Storage) DeleteFriend(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Or{
			sq.Eq{
				fieldFirstUserID:  authorID,
				fieldSecondUserID: recipientID,
				fieldDeletedAt:    nil,
			},
			sq.Eq{
				fieldFirstUserID:  recipientID,
				fieldSecondUserID: authorID,
				fieldDeletedAt:    nil,
			},
		}).
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

type friendEntity struct {
	ID         string `db:"id"`
	Username   string `db:"username"`
	FirstName  string `db:"first_name"`
	SecondName string `db:"second_name"`
}

func friendEntityToModel(entity friendEntity) *model.Friend {
	return &model.Friend{
		UserID:     model.UserID(entity.ID),
		Username:   entity.Username,
		FirstName:  entity.FirstName,
		SecondName: entity.SecondName,
	}
}

func friendEntitiesToModels(entities []friendEntity) []*model.Friend {
	var friendModels []*model.Friend
	for _, entity := range entities {
		friendModels = append(friendModels, friendEntityToModel(entity))
	}
	return friendModels
}
