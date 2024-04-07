package storage

import (
	"context"

	"social-network/internal/model"
)

type Storage interface {
	User() UserRepository
	Token() TokenRepository
}

type UserRepository interface {
	GetUserByLogin(ctx context.Context, userLogin *model.UserLogin) (*model.User, error)
	CreateUser(ctx context.Context, user *model.UserRegister) (*model.User, error)
	GetUserByID(ctx context.Context, id model.UserID) (*model.User, error)
	SearchUser(ctx context.Context, firstName, lastName string) ([]*model.User, error)
}

type TokenRepository interface {
	GetCurrentUserToken(ctx context.Context, id model.UserID) (*model.Token, error)
	CreateToken(ctx context.Context, token *model.Token) (*model.Token, error)
	RevokeToken(ctx context.Context, token *model.Token) error
	RefreshToken(ctx context.Context, token *model.Token) (*model.Token, error)
}
