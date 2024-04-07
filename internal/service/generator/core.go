package generator

import (
	"context"
	"fmt"

	"social-network/internal/model"
	"social-network/internal/storage"
)

type Service interface {
	BatchGenerateUsers(ctx context.Context, params *BatchGenerateUsersParams) error
}

type ServiceImpl struct {
	Storage storage.Storage
}

func (s *ServiceImpl) BatchGenerateUsers(ctx context.Context) error {
	//TODO: make atomic transaction and batch insert
	var usersList []*model.UserRegister
	for _, user := range usersList {
		_, err := s.Storage.User().CreateUser(ctx, user)
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}
	}

	return nil
}
