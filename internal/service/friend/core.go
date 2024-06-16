package friend

import (
	"context"
	"fmt"

	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/storage"
	"github.com/syth0le/social-network/internal/utils"
)

type Service interface {
	ListFriends(ctx context.Context, params *ListFriendsParams) ([]*model.Friend, error) // TODO: pagination
	AddFriend(ctx context.Context, params *AddFriendParams) error
	DeleteFriend(ctx context.Context, params *DeleteFriendParams) error
}

type ServiceImpl struct {
	Storage storage.Storage
}

type ListFriendsParams struct {
	UserID model.UserID
}

func (s *ServiceImpl) ListFriends(ctx context.Context, params *ListFriendsParams) ([]*model.Friend, error) {
	friends, err := s.Storage.Friend().ListFriends(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}

	return friends, nil
}

type AddFriendParams struct {
	AuthorID   model.UserID
	FollowerID model.UserID
}

func (s *ServiceImpl) AddFriend(ctx context.Context, params *AddFriendParams) error {
	err := s.Storage.Friend().AddFriend(ctx, &model.AddFriendParams{
		FID:        model.FriendID(utils.GenerateFUID()),
		SID:        model.FriendID(utils.GenerateFUID()),
		AuthorID:   params.AuthorID,
		FollowerID: params.FollowerID,
	})
	if err != nil {
		return fmt.Errorf("add friend: %w", err)
	}

	return nil
}

type DeleteFriendParams struct {
	AuthorID   model.UserID
	FollowerID model.UserID
}

func (s *ServiceImpl) DeleteFriend(ctx context.Context, params *DeleteFriendParams) error {
	err := s.Storage.Friend().DeleteFriend(ctx, params.AuthorID, params.FollowerID)
	if err != nil {
		return fmt.Errorf("delete friend: %w", err)
	}

	return nil
}
