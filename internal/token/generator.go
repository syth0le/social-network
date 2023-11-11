package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"social-network/internal/model"
)

type Generator struct {
	saltValue  string
	serverName string
}

func NewGenerator(saltValue, serverName string) *Generator {
	return &Generator{
		saltValue:  saltValue,
		serverName: serverName,
	}
}

func (g *Generator) GenerateToken(user *model.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss":      g.serverName,
		"sub":      user.Username,
		"sub_name": user.FirstName,
	})

	s, err := token.SignedString(g.saltValue)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}

	return s, nil
}
