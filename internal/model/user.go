package model

import (
	"fmt"
	"gopkg.in/validator.v2"
	"time"
)

type UserID string

func (id UserID) String() string {
	return string(id)
}

type TokenID string

func (id TokenID) String() string {
	return string(id)
}

type User struct {
	UserID     UserID
	Username   string
	FirstName  string
	SecondName string
	Sex        string
	Birthdate  time.Time
	Biography  string
	City       string
}

type Token struct {
	TokenID TokenID
	UserID  UserID
	Token   string
}

type TokenWithMetadata struct {
	TokenID  string    `validator:"nonzero"`
	UserID   UserID    `validator:"nonzero"`
	Token    string    `validator:"nonzero"`
	AlivedAt time.Time `validator:"nonzero"`
}

func (t *TokenWithMetadata) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

type UserLogin struct {
	Username       string
	HashedPassword string
}

type UserRegister struct {
	ID             string    `validator:"nonzero"`
	Username       string    `validator:"nonzero"`
	HashedPassword string    `validator:"nonzero"`
	FirstName      string    `validator:"nonzero"`
	SecondName     string    `validator:"nonzero"`
	Sex            string    `validator:"nonzero"`
	Birthdate      time.Time `validator:"nonzero"`
	Biography      string    `validator:"nonzero"`
	City           string    `validator:"nonzero"`
}

func (u *UserRegister) Validate() error {
	if err := validator.Validate(u); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
