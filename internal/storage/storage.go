package storage

import (
	"context"

	"social-network/internal/model"
)

type Storage interface {
	User() UserRepository
	Token() TokenRepository
	Friend() FriendRepository
	Post() PostRepository
}

type UserRepository interface {
	GetUserByLogin(ctx context.Context, userLogin *model.UserLogin) (*model.User, error)
	CreateUser(ctx context.Context, user *model.UserRegister) (*model.User, error)
	GetUserByID(ctx context.Context, id model.UserID) (*model.User, error)
	GetUserByTokenAndID(ctx context.Context, userID model.UserID, token string) (*model.User, error)
	SearchUser(ctx context.Context, firstName, lastName string) ([]*model.User, error)
	BatchCreateUser(ctx context.Context, user []*model.UserRegister) error
}

type TokenRepository interface {
	GetCurrentUserToken(ctx context.Context, id model.UserID) (*model.Token, error)
	CreateToken(ctx context.Context, token *model.Token) (*model.Token, error)
	RevokeToken(ctx context.Context, token *model.Token) error
	RefreshToken(ctx context.Context, token *model.Token) (*model.Token, error)
}

type FriendRepository interface {
	GetFriend(ctx context.Context, authorID, followerID model.UserID) (*model.Friend, error)
	ListFriends(ctx context.Context, userID model.UserID) ([]*model.Friend, error)
	AddFriend(ctx context.Context, params *model.AddFriendParams) error
	DeleteFriend(ctx context.Context, authorID, followerID model.UserID) error
}

type PostRepository interface {
	Create(ctx context.Context, params *model.CreatePostParams) error
	Get(ctx context.Context, postID model.PostID) (*model.Post, error)
	Update(ctx context.Context, postID model.PostID, text string) error
	Delete(ctx context.Context, postID model.PostID) error
	GetFeed(ctx context.Context, userID model.UserID) ([]*model.Post, error)
}
