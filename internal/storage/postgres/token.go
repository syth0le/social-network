package postgres

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"social-network/internal/model"
	"time"
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
		return nil, fmt.Errorf("incorrect sql") // todo internal error
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Slave(), &entity, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("internal error") // TODO: error wrapper SWITCH-CASE internal/not found
	}

	return tokenEntityToModel(entity), nil
}

func (s *Storage) CreateToken(ctx context.Context, params *model.TokenWithMetadata) (*model.Token, error) {
	err := params.Validate()
	if err != nil {
		return nil, fmt.Errorf("params validate: %w", err)
	}

	now := time.Now().Truncate(time.Millisecond)
	sql, args, err := sq.Insert(TokenTable).
		Columns(tokenFields...).
		Values(params.TokenID, params.UserID, params.Token, now, params.AlivedAt).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("incorrect sql") // todo internal error
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("internal error") // TODO: error wrapper SWITCH-CASE internal/not found
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
		return fmt.Errorf("incorrect sql") // todo internal error
	}

	result, err := s.Master().ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("internal error")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("internal error")
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not found error")
	}

	return nil
}

func (s *Storage) RefreshToken(ctx context.Context, params *model.TokenWithMetadata) (*model.Token, error) {
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
		Set(fieldAlivedAt, params.AlivedAt).
		Suffix(returningToken).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("incorrect sql") // todo internal error
	}

	var entity tokenEntity
	err = sqlx.GetContext(ctx, s.Master(), &entity, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("internal error") // TODO: error wrapper SWITCH-CASE internal/not found
	}

	return tokenEntityToModel(entity), nil
}

type tokenEntity struct {
	UserID string `db:"user_id"`
	Token  string `db:"token"`
}

func tokenEntityToModel(entity tokenEntity) *model.Token {
	return &model.Token{
		UserID: model.UserID(entity.UserID),
		Token:  entity.Token,
	}
}
