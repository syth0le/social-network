package queue

import (
	"context"
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"

	"social-network/internal/clients/rabbit"
	"social-network/internal/infrastructure_services/cache"
	"social-network/internal/model"
	"social-network/internal/service/friend"
)

type ConsumerService struct {
	Consumer      rabbit.Consumer
	Logger        *zap.Logger
	FriendService friend.Service
	CacheService  cache.Service
}

func (s *ConsumerService) CreateFeed(ctx context.Context, post *model.Post) error {
	// TODO: передавать значение по каналу из одного места в другое)
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

func (s *ConsumerService) Run() error {
	err := s.Consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("consumed: %v", string(d.Body))

		post := new(model.Post)
		err := post.UnmarshalBinary(d.Body)
		if err != nil {
			return rabbitmq.NackDiscard
		}

		err = s.CreateFeed(context.Background(), post)
		if err != nil {
			return rabbitmq.NackDiscard
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
	binary, err := post.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal binary: %w", err)
	}

	err = s.Producer.Publish(binary)
	if err != nil {
		return fmt.Errorf("producer publish: %w", err)
	}

	return nil
}
