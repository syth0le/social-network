package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/syth0le/gopnik/db/postgres"
	"go.uber.org/zap"

	"github.com/syth0le/social-network/internal/model"
)

type PostRepositoryTestSuite struct {
	suite.Suite
	ctx  context.Context
	repo *Storage
}

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}

func (s *PostRepositoryTestSuite) SetupSuite() {
	repo, err := NewStorage(zap.NewNop(), postgres.StorageConfig{})
	if err != nil {
		s.T().Fatal()
	}

	s.repo = repo
	s.ctx = context.Background()
}

func (s *PostRepositoryTestSuite) SetupTest() {
	s.repo.Delete(s.ctx, "IDlsfgldswjflsdfl")
}

func (s *PostRepositoryTestSuite) TearDownSuite() {
	s.repo.Delete(s.ctx, "IDlsfgldswjflsdfl")
}

func (s *PostRepositoryTestSuite) TestGetSuccess() {
	var tableData = map[string]struct {
		postID       model.PostID
		expectedPost model.Post
	}{
		"success test data 1": {"324234", model.Post{
			ID: "id", Text: "text", AuthorID: "ID",
		}},
		"success test data 2": {"222222", model.Post{
			ID: "id", Text: "text", AuthorID: "ID",
		}},
	}

	for name, data := range tableData {
		s.T().Run(name, func(t *testing.T) {
			result, err := s.repo.Get(s.ctx, data.postID)

			require.NoError(t, err)
			require.Equal(t, data.expectedPost, result)
		})
	}
}

func (s *PostRepositoryTestSuite) TestGetFailed() {
	var tableData = map[string]struct {
		postID model.PostID
	}{
		"failed data 1":    {""},
		"deleted post get": {"222222"},
	}

	for name, data := range tableData {
		s.T().Run(name, func(t *testing.T) {
			_, err := s.repo.Get(s.ctx, data.postID)

			require.Error(t, err)
		})
	}
}
