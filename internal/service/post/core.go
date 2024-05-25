package post

import (
	"context"
	"fmt"

	"social-network/internal/model"
	"social-network/internal/storage"
	"social-network/internal/utils"
)

type Service interface {
	Create(ctx context.Context, params *CreatePostParams) error
	GetPostByID(ctx context.Context, params *GetPostByIDParams) (*model.Post, error)
	Update(ctx context.Context, params *UpdatePostParams) error
	Delete(ctx context.Context, params *DeletePostParams) error
	GetFeed(ctx context.Context, params *GetFeedParams) ([]*model.Post, error)
}

type ServiceImpl struct {
	Storage storage.Storage
}

type CreatePostParams struct {
	UserID model.UserID
	Text   string
}

func (s ServiceImpl) Create(ctx context.Context, params *CreatePostParams) error {
	err := s.Storage.Post().Create(ctx, &model.CreatePostParams{
		ID:       model.PostID(utils.GeneratePUID()),
		Text:     params.Text,
		AuthorID: params.UserID,
	})
	if err != nil {
		return fmt.Errorf("create post: %w", err)
	}

	// вот тут ставить задачу в очередь (нужна транзакция)

	return nil
}

type GetPostByIDParams struct {
	PostID model.PostID
}

func (s ServiceImpl) GetPostByID(ctx context.Context, params *GetPostByIDParams) (*model.Post, error) {
	post, err := s.Storage.Post().Get(ctx, params.PostID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return post, nil
}

type UpdatePostParams struct {
	PostID model.PostID
	Text   string
}

func (s ServiceImpl) Update(ctx context.Context, params *UpdatePostParams) error {
	err := s.Storage.Post().Update(ctx, params.PostID, params.Text)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	return nil
}

type DeletePostParams struct {
	PostID model.PostID
}

func (s ServiceImpl) Delete(ctx context.Context, params *DeletePostParams) error {
	err := s.Storage.Post().Delete(ctx, params.PostID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

type GetFeedParams struct {
	FollowerID model.UserID
}

func (s ServiceImpl) GetFeed(ctx context.Context, params *GetFeedParams) ([]*model.Post, error) {
	// todo logic with queue and cache
	posts, err := s.Storage.Post().GetFeed(ctx, params.FollowerID)
	if err != nil {
		return nil, fmt.Errorf("get feed: %w", err)
	}

	return posts, nil
}
