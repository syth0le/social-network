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
	err = sqlx.GetContext(ctx, nil, &entity, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("internal error") // TODO: error wrapper SWITCH-CASE internal/not found
	}

	return tokenEntityToModel(entity), nil
}

func (s *Storage) CreateToken(ctx context.Context, token *model.Token) (*model.Token, error) {
	return nil, nil
}

func (s *Storage) RevokeToken(ctx context.Context, token *model.Token) (*model.Token, error) {
	return nil, nil
}

func (s *Storage) RefreshToken(ctx context.Context, token *model.Token) (*model.Token, error) {
	return nil, nil
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
