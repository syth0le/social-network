package postgres

import (
	"go.uber.org/zap"
	"social-network/cmd/social-network/configuration"
	"social-network/internal/storage"
)

type Storage struct {
	//storage *postgres.Storage
}

func NewStorage(logger *zap.Logger, config configuration.StorageConfig) (*Storage, error) {
	return &Storage{
		//storage: nil,
	}, nil
}

func (s *Storage) User() storage.UserRepository {
	return s
}

func (s *Storage) Token() storage.TokenRepository {
	return s
}
