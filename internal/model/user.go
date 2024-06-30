package model

import (
	"fmt"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"gopkg.in/validator.v2"
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
	UserID         UserID `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FirstName      string `json:"first_name"`
	SecondName     string `json:"second_name"`
	Sex            string `json:"sex"`
	Birthdate      time.Time
	Biography      string `json:"biography"`
	City           string `json:"city"`
}

type TarantoolUser struct {
	UserID         string `json:"id"`
	FirstName      string `json:"first_name"`
	SecondName     string `json:"second_name"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Sex            string `json:"sex"`
	Biography      string `json:"biography"`
	City           string `json:"city"`
}

func (t *TarantoolUser) DecodeMsgpack(d *msgpack.Decoder) error {
	var (
		err error
		l   int
	)
	if l, err = d.DecodeMapLen(); err != nil {
		return fmt.Errorf("decoding map length: %w", err)
	}
	for i := 0; i < l; i++ {
		key, err := d.DecodeInterface()
		if err != nil {
			return err
		}
		value, err := d.DecodeInterface()
		if err != nil {
			return err
		}
		switch key {
		case "id":
			t.UserID = value.(string)
		case "first_name":
			t.FirstName = value.(string)
		case "second_name":
			t.SecondName = value.(string)
		case "username":
			t.Username = value.(string)
		case "hashed_password":
			t.HashedPassword = value.(string)
		case "sex":
			t.Sex = value.(string)
		case "biography":
			t.Biography = value.(string)
		case "city":
			t.City = value.(string)
		}
	}
	return nil
}

type Token struct {
	TokenID        TokenID   `validator:"nonzero"`
	UserID         UserID    `validator:"nonzero"`
	Token          string    `validator:"nonzero"`
	ExpirationDate time.Time `validator:"nonzero"`
}

func (t *Token) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

type UserLogin struct {
	Username string
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
