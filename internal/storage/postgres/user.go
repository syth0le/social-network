package postgres

import (
	"context"
	"social-network/internal/model"
)

func (s *Storage) LoginUser(ctx context.Context, userLogin model.UserLogin) (*model.User, error) {
	return nil, nil
}

func (s *Storage) RegisterUser(ctx context.Context, user model.UserRegister) (*model.User, error) {
	return nil, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id model.UserID) (*model.User, error) {
	return nil, nil
}
