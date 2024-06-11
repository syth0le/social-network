package post

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"social-network/internal/infrastructure_services/cache"
	"social-network/internal/infrastructure_services/queue"
	"social-network/internal/model"
	"social-network/internal/storage"
	"social-network/internal/utils"
)

type Service interface {
	Create(ctx context.Context, params *CreatePostParams) (*model.Post, error)
	GetPostByID(ctx context.Context, params *GetPostByIDParams) (*model.Post, error)
	Update(ctx context.Context, params *UpdatePostParams) error
	Delete(ctx context.Context, params *DeletePostParams) error
	GetFeed(ctx context.Context, params *GetFeedParams) ([]*model.Post, error)
	GetLastPosts(ctx context.Context, params *GetLastPostsParams) ([]*model.Post, error)
}

type ServiceImpl struct {
	Storage         storage.Storage
	Cache           cache.Service
	Logger          *zap.Logger
	ProducerService *queue.ProducerService
}

type CreatePostParams struct {
	UserID model.UserID
	Text   string
}

func (s *ServiceImpl) Create(ctx context.Context, params *CreatePostParams) (*model.Post, error) {
	// TODO: make transaction???
	post, err := s.Storage.Post().Create(ctx, &model.CreatePostParams{
		ID:       model.PostID(utils.GeneratePUID()),
		Text:     params.Text,
		AuthorID: params.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	err = s.Cache.AddPost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("add post to cache: %w", err)
	}

	err = s.ProducerService.CreateFeed(post)
	if err != nil {
		return nil, fmt.Errorf("publish to queue: %w", err)
	}

	return post, nil
}

type GetPostByIDParams struct {
	PostID model.PostID
}

func (s *ServiceImpl) GetPostByID(ctx context.Context, params *GetPostByIDParams) (*model.Post, error) {
	post, err := s.Cache.GetPostByID(ctx, params.PostID)
	if err == nil {
		return post, nil
	}
	s.Logger.Sugar().Infof("cannot get post by id from cache: %s - %v", params.PostID, err)

	post, err = s.Storage.Post().Get(ctx, params.PostID)
	if err != nil {
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	return post, nil
}

type UpdatePostParams struct {
	PostID model.PostID
	Text   string
}

func (s *ServiceImpl) Update(ctx context.Context, params *UpdatePostParams) error {
	err := s.Storage.Post().Update(ctx, params.PostID, params.Text)
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	// TODO: invalidate cache

	return nil
}

type DeletePostParams struct {
	PostID model.PostID
}

func (s *ServiceImpl) Delete(ctx context.Context, params *DeletePostParams) error {
	err := s.Storage.Post().Delete(ctx, params.PostID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	// TODO: invalidate cache from feeds (hard logic)
	err = s.Cache.DeletePost(ctx, params.PostID)
	if err != nil {
		s.Logger.Sugar().Warnf("cannot delete from cache: %v", err)
		return nil
	}

	return nil
}

type GetFeedParams struct {
	FollowerID model.UserID
}

func (s *ServiceImpl) GetFeed(ctx context.Context, params *GetFeedParams) ([]*model.Post, error) {
	feed, err := s.Cache.GetFeedByUserID(ctx, params.FollowerID)
	if err == nil {
		return feed, nil
	}

	s.Logger.Sugar().Infof("cannot get feed by user id from cache: %s - %v", params.FollowerID, err)

	posts, err := s.Storage.Post().GetFeed(ctx, params.FollowerID)
	if err != nil {
		return nil, fmt.Errorf("get feed: %w", err)
	}

	return posts, nil
}

type GetLastPostsParams struct {
	Duration time.Duration
}

func (s *ServiceImpl) GetLastPosts(ctx context.Context, params *GetLastPostsParams) ([]*model.Post, error) {
	if params.Duration == 0 {
		return nil, fmt.Errorf("duration cannot be nil")
	}

	posts, err := s.Storage.Post().GetLastPosts(ctx, params.Duration)
	if err != nil {
		return nil, fmt.Errorf("get feed: %w", err)
	}

	return posts, nil
}
