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

func (s *Storage) Create(ctx context.Context, params *model.CreatePostParams) (*model.Post, error) {
	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("params validate: %w", err)
	}

	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Insert(PostTable).
		Columns(postFields...).
		Values(params.ID, params.AuthorID, params.Text, now, now).
		Suffix(returningPost).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity postEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return postEntityToModel(entity), nil
}

func (s *Storage) Get(ctx context.Context, postID model.PostID) (*model.Post, error) {
	sql, args, err := sq.Select(postFields...).
		From(PostTable).
		Where(sq.Eq{
			fieldID:        postID,
			fieldDeletedAt: nil,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity postEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return postEntityToModel(entity), nil
}

func (s *Storage) Update(ctx context.Context, postID model.PostID, text string) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(PostTable).
		Set(fieldText, text).
		Set(fieldUpdatedAt, now).
		Where(sq.Eq{
			fieldID:        postID,
			fieldDeletedAt: nil,
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

func (s *Storage) Delete(ctx context.Context, postID model.PostID) error {
	now := time.Now().Truncate(time.Millisecond)

	sql, args, err := sq.Update(PostTable).
		Set(fieldDeletedAt, now).
		Where(sq.Eq{
			fieldID:        postID,
			fieldDeletedAt: nil,
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

func (s *Storage) GetFeed(ctx context.Context, userID model.UserID) ([]*model.Post, error) {
	sql, args, err := sq.Select(
		tableFields(PostTable, postFields)...,
	).From(PostTable).
		Join(joinString(PostTable, fieldUserID, UserTable, fieldID)).
		Join(joinString(PostTable, fieldUserID, FriendTable, fieldFirstUserID)).
		Where(sq.Eq{
			tableField(PostTable, fieldDeletedAt):      nil,
			tableField(FriendTable, fieldSecondUserID): userID,
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entities []postEntity
	err = sqlx.SelectContext(ctx, s.Slave(), &entities, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return postEntitiesToModels(entities), nil
}

func (s *Storage) GetLastPosts(ctx context.Context, duration time.Duration) ([]*model.Post, error) {
	now := time.Now().Truncate(time.Millisecond)
	creationDateThreshold := now.Add(-duration)

	sql, args, err := sq.Select(postFields...).
		From(PostTable).
		Where(
			sq.And{
				sq.GtOrEq{
					fieldCreatedAt: creationDateThreshold,
				},
				sq.Eq{
					fieldDeletedAt: nil,
				},
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, xerrors.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entities []postEntity
	err = sqlx.SelectContext(ctx, s.Slave(), &entities, sql, args...)
	if err != nil {
		return nil, xerrors.WrapSqlError(err)
	}

	return postEntitiesToModels(entities), nil
}

type postEntity struct {
	ID        string `db:"id"`
	Text      string `db:"text"`
	AuthorID  string `db:"user_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func postEntityToModel(entity postEntity) *model.Post {
	return &model.Post{
		ID:       model.PostID(entity.ID),
		Text:     entity.Text,
		AuthorID: model.UserID(entity.AuthorID),
	}
}

func postEntitiesToModels(entities []postEntity) []*model.Post {
	var postModels []*model.Post
	for _, entity := range entities {
		postModels = append(postModels, postEntityToModel(entity))
	}
	return postModels
}
