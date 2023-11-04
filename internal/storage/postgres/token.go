package postgres

import (
	"context"
	"social-network/internal/model"
)

func (s *Storage) GetCurrentUserToken(ctx context.Context, id model.UserID) (*model.Token, error) {
	return nil, nil
}

func (s *Storage) CreateToken(ctx context.Context, token model.Token) (*model.Token, error) {
	return nil, nil
}

func (s *Storage) RevokeToken(ctx context.Context, token model.Token) (*model.Token, error) {
	return nil, nil
}

func (s *Storage) RefreshToken(ctx context.Context, token model.Token) (*model.Token, error) {
	return nil, nil
}
