package heater

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"social-network/internal/infrastructure_services/cache"
	"social-network/internal/service/friend"
	"social-network/internal/service/post"
)

type HeaterService struct {
	Logger         *zap.Logger
	CacheService   cache.Service
	PostService    post.Service
	FriendService  friend.Service
	HeaterDuration time.Duration
}

func (s *HeaterService) Run(ctx context.Context) error {
	s.Logger.Info("heating cache started....")

	posts, err := s.PostService.GetLastPosts(ctx, &post.GetLastPostsParams{
		Duration: s.HeaterDuration,
	})
	if err != nil {
		return fmt.Errorf("get last posts: %w", err)
	}

	// не проверяем пустой ли кэш, потому что даже если запишем значения - то что было там перетрется
	for _, postEntity := range posts {
		friendsList, err := s.FriendService.ListFriends(ctx, &friend.ListFriendsParams{UserID: postEntity.AuthorID})
		if err != nil {
			return fmt.Errorf("list friends: %w", err)
		}

		for _, friendEntity := range friendsList {
			err := s.CacheService.AddPostForUser(ctx, friendEntity.UserID, postEntity) // TODO: make batch create
			if err != nil {
				return fmt.Errorf("add post for user: %w", err)
			}
		}

		err = s.CacheService.AddPost(ctx, postEntity)
		if err != nil {
			return fmt.Errorf("add post: %w", err)
		}
	}

	s.Logger.Info("heating cache successfully finished")

	return nil
}
