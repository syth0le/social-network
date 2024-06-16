package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"github.com/syth0le/social-network/internal/clients/rabbit"
	"github.com/syth0le/social-network/internal/infrastructure_services/cache"
	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/service/friend"
)

type ConsumerService struct {
	Consumer      rabbit.Consumer
	Logger        *zap.Logger
	FriendService friend.Service
	CacheService  cache.Service
}

func (s *ConsumerService) CreateFeed(ctx context.Context, post *model.Post) error {
	friends, err := s.FriendService.ListFriends(ctx, &friend.ListFriendsParams{UserID: post.AuthorID})
	if err != nil {
		return fmt.Errorf("list friends: %w", err)
	}

	for _, user := range friends {
		err := s.CacheService.AddPostForUser(ctx, user.UserID, post) // TODO: make create tasks for ranges (1-10)
		if err != nil {
			return fmt.Errorf("add post for user to cache: %w", err)
		}
	}

	return nil
}

func (s *ConsumerService) UpdatePostInFeed(ctx context.Context, previousPost, post *model.Post) error {
	friends, err := s.FriendService.ListFriends(ctx, &friend.ListFriendsParams{UserID: post.AuthorID})
	if err != nil {
		return fmt.Errorf("list friends: %w", err)
	}

	s.Logger.Sugar().Infof("friends: %#v", friends)

	for _, user := range friends {
		// TODO: make transaction
		err := s.CacheService.DeletePostForUser(ctx, user.UserID, previousPost)
		if err != nil {
			return fmt.Errorf("delete post for user from cache: %w", err)
		}

		err = s.CacheService.AddPostForUser(ctx, user.UserID, post)
		if err != nil {
			return fmt.Errorf("add post for user to cache: %w", err)
		}
	}

	return nil
}

func (s *ConsumerService) DeletePostFromFeed(ctx context.Context, post *model.Post) error {
	friends, err := s.FriendService.ListFriends(ctx, &friend.ListFriendsParams{UserID: post.AuthorID})
	if err != nil {
		return fmt.Errorf("list friends: %w", err)
	}

	for _, user := range friends {
		err := s.CacheService.DeletePostForUser(ctx, user.UserID, post)
		if err != nil {
			return fmt.Errorf("delete post for user from cache: %w", err)
		}
	}

	return nil
}

func (s *ConsumerService) Run() error {
	err := s.Consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))

		post := new(model.PostAction)
		err := post.UnmarshalBinary(d.Body)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		switch post.Action {
		case model.CreateAction:
			err = s.CreateFeed(context.Background(), post.Post)
			if err != nil {
				return rabbitmq.NackDiscard
			}
		case model.DeleteAction:
			err = s.DeletePostFromFeed(context.Background(), post.Post)
			if err != nil {
				return rabbitmq.NackDiscard
			}
		}

		// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
		return rabbitmq.Ack
	})
	if err != nil {
		return fmt.Errorf("consumer run: %w", err)
	}

	return nil
}

type ProducerService struct {
	Producer rabbit.Publisher
	Logger   *zap.Logger
}

func (s *ProducerService) CreateFeed(post *model.Post) error {
	return s.makeFeed(post, model.CreateAction)
}

func (s *ProducerService) DeleteFromFeed(post *model.Post) error {
	return s.makeFeed(post, model.DeleteAction)
}

func (s *ProducerService) makeFeed(post *model.Post, action model.Action) error {
	postAction := model.PostAction{
		Action: action,
		Post:   post,
	}
	binary, err := postAction.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal binary: %w", err)
	}

	err = s.Producer.Publish(binary)
	if err != nil {
		return fmt.Errorf("producer publish: %w", err)
	}

	return nil
}
