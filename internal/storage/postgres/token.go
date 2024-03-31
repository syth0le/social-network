package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"social-network/internal/model"
	"social-network/internal/utils"
)

func (s *Storage) GetCurrentUserToken(ctx context.Context, id model.UserID) (*model.Token, error) {
	sql, args, err := sq.Select(tokenFields...).
		From(TokenTable).
		Where(
			sq.Eq{
				fieldUserID:    id.String(),
				fieldDeletedAt: nil,
			},
			sq.LtOrEq{
				fieldAlivedAt: time.Now(),
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return tokenEntityToModel(entity), nil
}

func (s *Storage) CreateToken(ctx context.Context, params *model.Token) (*model.Token, error) {
	err := params.Validate()
	if err != nil {
		return nil, utils.WrapValidationError(fmt.Errorf("params validate: %w", err))
	}

	now := time.Now().Truncate(time.Millisecond)
	sql, args, err := sq.Insert(TokenTable).
		Columns(tokenFields...).
		Values(params.TokenID, params.UserID, params.Token, now, params.ExpirationDate).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return tokenEntityToModel(entity), nil
}

func (s *Storage) RevokeToken(ctx context.Context, params *model.Token) error {
	now := time.Now().Truncate(time.Millisecond)
	sql, args, err := sq.Update(TokenTable).
		Where(sq.Eq{
			fieldToken:     params.Token,
			fieldUserID:    params.UserID.String(),
			fieldDeletedAt: nil,
		}).
		Set(fieldDeletedAt, now).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	result, err := s.Master().ExecContext(ctx, sql, args...)
	if err != nil {
		return utils.WrapSqlError(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.WrapInternalError(err)
	}
	if rowsAffected == 0 {
		return utils.WrapNotFoundError(err, utils.NotFoundMessage)
	}

	return nil
}

func (s *Storage) RefreshToken(ctx context.Context, params *model.Token) (*model.Token, error) {
	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("params validate: %w", err)
	}
	sql, args, err := sq.Update(TokenTable).
		Where(sq.Eq{
			fieldToken:     params.Token,
			fieldUserID:    params.UserID.String(),
			fieldDeletedAt: nil,
		}).
		Set(fieldAlivedAt, params.ExpirationDate).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, utils.WrapInternalError(fmt.Errorf("incorrect sql"))
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, utils.WrapSqlError(err)
	}

	return tokenEntityToModel(entity), nil
}

type tokenEntity struct {
	ID        string    `db:"id"`
	UserID    string    `db:"user_id"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	AlivedAt  time.Time `db:"alived_at"`
}

func tokenEntityToModel(entity tokenEntity) *model.Token {
	return &model.Token{
		TokenID:        model.TokenID(entity.ID),
		UserID:         model.UserID(entity.UserID),
		Token:          entity.Token,
		ExpirationDate: entity.AlivedAt,
	}
}
