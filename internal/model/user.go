package model

type UserID string

func (id UserID) String() string {
	return string(id)
}

type User struct {
	UserID     UserID
	Username   string
	FirstName  string
	SecondName string
	Age        int
	Sex        string
	Birthdate  string
	Biography  string
	City       string
}

type Token struct {
	UserID UserID
	Token  string
}

type UserLogin struct {
	Username       string
	HashedPassword string
}

type UserRegister struct {
	Username       string
	HashedPassword string
	FirstName      string
	SecondName     string
	Age            int
	Sex            string
	Birthdate      string
	Biography      string
	City           string
}
