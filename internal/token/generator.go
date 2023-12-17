package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"social-network/cmd/social-network/configuration"
	"social-network/internal/model"
	"time"
)

const ExpirationDuration = time.Hour * 24

type Generator struct {
	saltValue  string
	serverName string
}

func NewGenerator(config configuration.ApplicationConfig) *Generator {
	return &Generator{
		saltValue:  config.SaltValue,
		serverName: config.App,
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

func (g *Generator) GetExpirationDate() time.Time {
	return time.Now().Add(ExpirationDuration)
}
