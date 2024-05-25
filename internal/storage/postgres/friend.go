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

//SELECT
//u.id, u.username, u.first_name, u.second_name
//FROM friend_table AS f
//JOIN user_table AS u ON
//CASE
//WHEN f.first_user_id = ID THEN f.first_user_id = u.id
//ELSE f.second_user_id = u.id
//END
//WHERE (f.first_user_id = ID OR f.second_user_id = ID) AND f.is_friend=true;

func (s *Storage) ListFriends(ctx context.Context, userID model.UserID) ([]*model.Friend, error) {
	sql, args, err := sq.Select(
		tableField(UserTable, fieldID),
		tableField(UserTable, fieldUsername),
		tableField(UserTable, fieldFirstName),
		tableField(UserTable, fieldSecondName),
	).
		Join(
			sq.Case().When(
				sq.Eq{tableField(FriendTable, fieldFirstUserID): userID},
				joinString(FriendTable, fieldFirstUserID, UserTable, fieldID), // TODO join case
			).Else(joinString(FriendTable, fieldSecondUserID, UserTable, fieldID)),
		).
		Where(sq.And{
			sq.Or{
				sq.Eq{tableField(FriendTable, fieldFirstUserID): userID},
				sq.Eq{tableField(FriendTable, fieldSecondUserID): userID},
			},
			sq.Eq{tableField(FriendTable, fieldStatus): fieldStatusAccepted},
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

// ListFollowers -- get followers
// SELECT * FROM friend_table AS f
// JOIN user_table AS u ON f.first_user_id = u.id
// WHERE f.first_user_id = ID AND f.is_friend=false;
func (s *Storage) ListFollowers(ctx context.Context, userID model.UserID) ([]*model.Friend, error) {
	sql, args, err := sq.Select(
		tableField(UserTable, fieldID),
		tableField(UserTable, fieldUsername),
		tableField(UserTable, fieldFirstName),
		tableField(UserTable, fieldSecondName),
	).
		Join(joinString(FriendTable, fieldFirstUserID, UserTable, fieldID)). // TODO неверный Join()
		Where(
			sq.Or{
				sq.Eq{
					tableField(FriendTable, fieldSecondUserID): userID,
					tableField(FriendTable, fieldStatus):       []string{fieldStatusExpected, fieldStatusDeclined},
				},
				sq.Eq{
					tableField(FriendTable, fieldFirstUserID): userID,
					tableField(FriendTable, fieldStatus):      fieldStatusRevoked,
				},
			},
		).
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

// ListSubscriptions -- get followed
// SELECT * FROM friend_table AS f
// JOIN user_table AS u ON f.second_user_id = u.id
// WHERE f.second_user_id = ID AND f.is_friend=false;
func (s *Storage) ListSubscriptions(ctx context.Context, userID model.UserID) ([]*model.Friend, error) {
	sql, args, err := sq.Select(
		tableField(UserTable, fieldID),
		tableField(UserTable, fieldUsername),
		tableField(UserTable, fieldFirstName),
		tableField(UserTable, fieldSecondName),
	).
		Join(joinString(FriendTable, fieldSecondUserID, UserTable, fieldID)). // TODO неверный Join()
		Where(
			sq.Or{
				sq.Eq{
					tableField(FriendTable, fieldFirstUserID): userID,
					tableField(FriendTable, fieldStatus):      []string{fieldStatusExpected, fieldStatusDeclined},
				},
				sq.Eq{
					tableField(FriendTable, fieldSecondUserID): userID,
					tableField(FriendTable, fieldStatus):       fieldStatusRevoked,
				},
			},
		).
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

func (s *Storage) SetFriendRequest(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldFirstUserID:  authorID,
			fieldSecondUserID: recipientID,
			fieldDeletedAt:    nil,
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

func (s *Storage) ConfirmFriendRequest(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldFirstUserID:  authorID,
			fieldSecondUserID: recipientID,
			fieldDeletedAt:    nil,
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

func (s *Storage) DeclineFriendRequest(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldFirstUserID:  authorID,
			fieldSecondUserID: recipientID,
			fieldDeletedAt:    nil,
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

func (s *Storage) RevokeFriendRequest(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldFirstUserID:  authorID,
			fieldSecondUserID: recipientID,
			fieldDeletedAt:    nil,
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

func (s *Storage) DeleteFriend(ctx context.Context, authorID, recipientID model.UserID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(FriendTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldFirstUserID:  authorID,
			fieldSecondUserID: recipientID,
			fieldDeletedAt:    nil,
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
