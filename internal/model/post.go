package model

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type PostID string

func (id PostID) String() string {
	return string(id)
}

type Post struct {
	ID       PostID `validator:"nonzero"`
	Text     string `validator:"nonzero"`
	Author   string `validator:"nonzero"`
	AuthorID UserID `validator:"nonzero"`
}

func (t *Post) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

type CreatePostParams struct {
	ID       PostID `validator:"nonzero"`
	Text     string `validator:"nonzero"`
	AuthorID UserID `validator:"nonzero"`
}

func (t *CreatePostParams) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
