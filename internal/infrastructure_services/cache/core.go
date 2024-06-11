package cache

import (
	"context"
	"fmt"

	xerrors "github.com/syth0le/gopnik/errors"
	"go.uber.org/zap"

	"social-network/internal/clients/redis"
	"social-network/internal/model"
)

const (
	PostHashType HashType = "post"
	UserHashType HashType = "user"

	defaultRange = 1000 // TODO: pagination
)

type HashType string

func (h HashType) String() string {
	return string(h)
}

type Service interface {
	AddPost(ctx context.Context, post *model.Post) error
	DeletePost(ctx context.Context, id model.PostID) error
	GetPostByID(ctx context.Context, id model.PostID) (*model.Post, error)
	AddPostForUser(ctx context.Context, userID model.UserID, post *model.Post) error
	GetFeedByUserID(ctx context.Context, id model.UserID) ([]*model.Post, error)
	DeletePostForUser(ctx context.Context, userID model.UserID, post *model.Post) error
}

type ServiceImpl struct {
	Client redis.Client
	Logger *zap.Logger
}

func (s *ServiceImpl) AddPost(ctx context.Context, post *model.Post) error {
	keyHash, err := makeHash(PostHashType, post.ID.String())
	if err != nil {
		return fmt.Errorf("make hash: %w", err)
	}

	err = s.Client.HSetNX(
		ctx,
		true,
		keyHash,
		PostHashType.String(),
		post,
	)
	if err != nil {
		return fmt.Errorf("cache set: %w", err)
	}

	s.Logger.Sugar().Infof("key %s saved in cache", keyHash)

	return nil
}

func (s *ServiceImpl) DeletePost(ctx context.Context, id model.PostID) error {
	keyHash, err := makeHash(PostHashType, id.String())
	if err != nil {
		return fmt.Errorf("make hash: %w", err)
	}

	err = s.Client.Delete(ctx, keyHash)
	if err != nil {
		return fmt.Errorf("cache delete: %w", err)
	}

	return nil
}

func (s *ServiceImpl) AddPostForUser(ctx context.Context, userID model.UserID, post *model.Post) error {
	keyHash, err := makeHash(UserHashType, userID.String())
	if err != nil {
		return fmt.Errorf("make hash: %w", err)
	}

	err = s.Client.LPush(ctx, keyHash, post)
	if err != nil {
		return fmt.Errorf("cache lpush: %w", err)
	}

	s.Logger.Sugar().Debugf("key %s saved in cache", keyHash)

	return nil
}

func (s *ServiceImpl) DeletePostForUser(ctx context.Context, userID model.UserID, post *model.Post) error {
	keyHash, err := makeHash(UserHashType, userID.String())
	if err != nil {
		return fmt.Errorf("make hash: %w", err)
	}

	err = s.Client.LRem(ctx, keyHash, post)
	if err != nil {
		return fmt.Errorf("cache lpush: %w", err)
	}

	s.Logger.Sugar().Debugf("key %s saved in cache", keyHash)

	return nil
}

func (s *ServiceImpl) GetPostByID(ctx context.Context, id model.PostID) (*model.Post, error) {
	keyHash, err := makeHash(PostHashType, id.String())
	if err != nil {
		return nil, fmt.Errorf("make hash: %w", err)
	}

	post := &model.Post{}
	err = s.Client.HGet(ctx, keyHash, PostHashType.String(), post)
	if err != nil {
		return nil, fmt.Errorf("key %s not found in Redis cache: %w", id.String(), err)
	}

	return post, nil
}

func (s *ServiceImpl) GetFeedByUserID(ctx context.Context, id model.UserID) ([]*model.Post, error) {
	keyHash, err := makeHash(UserHashType, id.String())
	if err != nil {
		return nil, fmt.Errorf("make hash: %w", err)
	}

	listPosts, err := s.Client.LRange(ctx, keyHash, 0, defaultRange) // TODO: make pagination
	if err != nil {
		return nil, fmt.Errorf("key %s not found in Redis cache: %w", id.String(), err)
	}

	posts := make([]*model.Post, len(listPosts))
	acc := 0
	for _, val := range listPosts {
		post := new(model.Post)
		if err = post.UnmarshalBinary([]byte(val)); err != nil {
			return nil, fmt.Errorf("unmarshal binary: %w", err)
		}

		posts[acc] = post
		acc += 1
	}

	return posts, nil
}

func makeHash(hashType HashType, key string) (string, error) {
	switch hashType {
	case PostHashType, UserHashType:
	default:
		return "", xerrors.WrapInternalError(fmt.Errorf("unexpected hash type: %s", hashType))
	}

	if key == "" {
		return "", xerrors.WrapInternalError(fmt.Errorf("key cannot be empty"))
	}

	return fmt.Sprintf("%s-%s", hashType, key), nil
}
