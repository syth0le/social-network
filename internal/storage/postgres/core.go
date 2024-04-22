package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	xstorage "github.com/syth0le/gopnik/db/postgres"
	"go.uber.org/zap"

	"social-network/internal/storage"
)

type Storage struct {
	storage *xstorage.PGStorage
}

func NewStorage(logger *zap.Logger, config xstorage.StorageConfig) (*Storage, error) {
	postgresStorage, err := xstorage.NewPGStorage(logger, config)
	if err != nil {
		return nil, fmt.Errorf("new pg storage: %w", err)
	}
	return &Storage{
		storage: postgresStorage,
	}, nil
}

func (s *Storage) User() storage.UserRepository {
	return s
}

func (s *Storage) Token() storage.TokenRepository {
	return s
}

func (s *Storage) Friend() storage.FriendRepository {
	return s
}

func (s *Storage) Post() storage.PostRepository {
	return s
}

func (s *Storage) Close() error {
	return s.storage.Close()
}

func (s *Storage) Master() sqlx.ExtContext {
	return s.storage.Master()
}

func (s *Storage) Slave() sqlx.ExtContext {
	return s.storage.Slave()
}

func (s *Storage) now() {
	//TODO: implement
}
