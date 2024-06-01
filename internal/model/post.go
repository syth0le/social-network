package model

import (
	"encoding/json"
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
	AuthorID UserID `validator:"nonzero"`
}

func (p *Post) Validate() error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func (p *Post) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Post) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}

	return nil
}

type CreatePostParams struct {
	ID       PostID `validator:"nonzero"`
	Text     string `validator:"nonzero"`
	AuthorID UserID `validator:"nonzero"`
}

func (p *CreatePostParams) Validate() error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
