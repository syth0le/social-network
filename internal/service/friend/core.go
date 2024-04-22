package friend

import (
	"context"
	"fmt"

	"social-network/internal/model"
	"social-network/internal/storage"
)

type Service interface {
	ListFriends(ctx context.Context, params *ListFriendsParams) ([]*model.Friend, error)             // TODO: pagination
	ListFollowers(ctx context.Context, params *ListFollowersParams) ([]*model.Friend, error)         // TODO: pagination
	ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) ([]*model.Friend, error) // TODO: pagination
	SetFriend(ctx context.Context, params *SetFriendParams) error
	DeleteFriend(ctx context.Context, params *DeleteFriendParams) error
}

type ServiceImpl struct {
	Storage storage.Storage
}

type ListFriendsParams struct {
	UserID model.UserID
}

func (s ServiceImpl) ListFriends(ctx context.Context, params *ListFriendsParams) ([]*model.Friend, error) {
	friends, err := s.Storage.Friend().ListFriends(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("list friends: %w", err)
	}

	return friends, nil
}

type ListFollowersParams struct {
	UserID model.UserID
}

func (s ServiceImpl) ListFollowers(ctx context.Context, params *ListFollowersParams) ([]*model.Friend, error) {
	followers, err := s.Storage.Friend().ListFollowers(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("list followers: %w", err)
	}

	return followers, nil
}

type ListSubscriptionsParams struct {
	UserID model.UserID
}

func (s ServiceImpl) ListSubscriptions(ctx context.Context, params *ListSubscriptionsParams) ([]*model.Friend, error) {
	subscriptions, err := s.Storage.Friend().ListSubscriptions(ctx, params.UserID)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}

	return subscriptions, nil
}

type SetFriendParams struct {
	AuthorID   model.UserID
	FollowerID model.UserID
}

func (s ServiceImpl) SetFriend(ctx context.Context, params *SetFriendParams) error {
	err := s.Storage.Friend().SetFriend(ctx, params.AuthorID, params.FollowerID)
	if err != nil {
		return fmt.Errorf("set friend: %w", err)
	}

	return nil
}

type DeleteFriendParams struct {
	AuthorID   model.UserID
	FollowerID model.UserID
}

func (s ServiceImpl) DeleteFriend(ctx context.Context, params *DeleteFriendParams) error {
	err := s.Storage.Friend().DeleteFriend(ctx, params.AuthorID, params.FollowerID)
	if err != nil {
		return fmt.Errorf("delete friend: %w", err)
	}

	return nil
}
