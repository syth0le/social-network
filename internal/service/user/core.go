package user

import (
	"context"
	"fmt"
	"time"

	"social-network/internal/model"
	"social-network/internal/storage"
	"social-network/internal/token"
	"social-network/internal/utils"
)

type Service interface {
	Login(ctx context.Context, params *LoginParams) (*model.Token, error)
	Register(ctx context.Context, params *RegisterParams) (*model.Token, error)
	GetUserByID(ctx context.Context, params *GetUserByIDParams) (*model.User, error)
	SearchUser(ctx context.Context, params *SearchUserParams) (*model.User, error)
}

type ServiceImpl struct {
	Storage        storage.Storage
	TokenGenerator *token.Generator
}

type LoginParams struct {
	Username string
	Password string
}

func (s *ServiceImpl) Login(ctx context.Context, params *LoginParams) (*model.Token, error) {
	//TODO: make atomic transaction
	user, err := s.Storage.User().GetUserByLogin(ctx, &model.UserLogin{
		Username: params.Username,
	})
	if err != nil {
		return nil, fmt.Errorf("login user: %w", err)
	}

	if err = utils.CheckPasswordHash(user.HashedPassword, params.Password); err != nil {
		return nil, utils.WrapNotFoundError(fmt.Errorf("not correct password: %w", err), "not correct credentials")
	}

	generatedToken, err := s.TokenGenerator.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	token, err := s.Storage.Token().CreateToken(ctx, generatedToken)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	return token, nil
}

type RegisterParams struct {
	Username   string
	Password   string
	FirstName  string
	SecondName string
	Age        int
	Sex        string
	Birthdate  time.Time
	Biography  string
	City       string
}

func (s *ServiceImpl) Register(ctx context.Context, params *RegisterParams) (*model.Token, error) {
	hashedPassword, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.Storage.User().CreateUser(ctx, &model.UserRegister{
		ID:             utils.GenerateUUID(),
		Username:       params.Username,
		HashedPassword: hashedPassword,
		FirstName:      params.FirstName,
		SecondName:     params.SecondName,
		Sex:            params.Sex,
		Birthdate:      params.Birthdate,
		Biography:      params.Biography,
		City:           params.City,
	})
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	generatedToken, err := s.TokenGenerator.GenerateToken(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	token, err := s.Storage.Token().CreateToken(ctx, generatedToken)
	if err != nil {
		return nil, fmt.Errorf("create token: %w", err)
	}

	return token, nil
}

type GetUserByIDParams struct {
	UserID model.UserID
}

func (s *ServiceImpl) GetUserByID(ctx context.Context, params *GetUserByIDParams) (*model.User, error) {
	user, err := s.Storage.User().GetUserByID(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

type SearchUserParams struct {
	FirstName string
	LastName  string
}

func (s *ServiceImpl) SearchUser(ctx context.Context, params *SearchUserParams) (*model.User, error) {
	user, err := s.Storage.User().SearchUser(ctx, params.FirstName, params.LastName)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}
