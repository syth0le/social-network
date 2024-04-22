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
	SetFriendRequest(ctx context.Context, params *SetFriendRequestParams) error
	ConfirmFriendRequest(ctx context.Context, params *ConfirmFriendRequestParams) error
	DeclineFriendRequest(ctx context.Context, params *DeclineFriendRequestParams) error
	RevokeFriendRequest(ctx context.Context, params *RevokeFriendRequestParams) error
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

type SetFriendRequestParams struct {
	AuthorID    model.UserID
	RecipientID model.UserID
}

func (s ServiceImpl) SetFriendRequest(ctx context.Context, params *SetFriendRequestParams) error {
	err := s.Storage.Friend().SetFriendRequest(ctx, params.AuthorID, params.RecipientID)
	if err != nil {
		return fmt.Errorf("set friend: %w", err)
	}

	return nil
}

type ConfirmFriendRequestParams struct {
	AuthorID    model.UserID
	RecipientID model.UserID
}

func (s ServiceImpl) ConfirmFriendRequest(ctx context.Context, params *ConfirmFriendRequestParams) error {
	err := s.Storage.Friend().ConfirmFriendRequest(ctx, params.AuthorID, params.RecipientID)
	if err != nil {
		return fmt.Errorf("confirm friend request: %w", err)
	}

	return nil
}

type DeclineFriendRequestParams struct {
	AuthorID    model.UserID
	RecipientID model.UserID
}

func (s ServiceImpl) DeclineFriendRequest(ctx context.Context, params *DeclineFriendRequestParams) error {
	err := s.Storage.Friend().DeclineFriendRequest(ctx, params.AuthorID, params.RecipientID)
	if err != nil {
		return fmt.Errorf("decline friend request: %w", err)
	}

	return nil
}

type RevokeFriendRequestParams struct {
	AuthorID    model.UserID
	RecipientID model.UserID
}

func (s ServiceImpl) RevokeFriendRequest(ctx context.Context, params *RevokeFriendRequestParams) error {
	err := s.Storage.Friend().RevokeFriendRequest(ctx, params.AuthorID, params.RecipientID)
	if err != nil {
		return fmt.Errorf("revoke friend request: %w", err)
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
