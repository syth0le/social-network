package model

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type FriendID string

func (id FriendID) String() string {
	return string(id)
}

type Friend struct {
	UserID     UserID
	Username   string
	FirstName  string
	SecondName string
}

type AddFriendParams struct {
	FID        FriendID `validator:"nonzero"`
	SID        FriendID `validator:"nonzero"`
	AuthorID   UserID   `validator:"nonzero"`
	FollowerID UserID   `validator:"nonzero"`
}

func (t *AddFriendParams) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
