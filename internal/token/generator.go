package token

import (
	"fmt"
	"time"

	"social-network/cmd/social-network/configuration"
	"social-network/internal/model"
	"social-network/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

const ExpirationDuration = time.Hour * 24 // todo: get expiration time from settings

type Generator struct {
	saltValue  []byte
	serverName string
}

func NewGenerator(config configuration.ApplicationConfig) *Generator {
	return &Generator{
		saltValue:  []byte(config.SaltValue),
		serverName: config.App,
	}
}

func (g *Generator) GenerateToken(user *model.User) (*model.Token, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be empty")
	}

	expirationDate := g.getExpirationDate()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":        g.serverName,
		"sub":        user.Username,
		"sub_name":   user.FirstName,
		"expires_at": expirationDate.Unix(),
	})

	tokenString, err := token.SignedString(g.saltValue)
	if err != nil {
		return nil, fmt.Errorf("signed string: %w", err)
	}

	return &model.Token{
		TokenID:        model.TokenID(utils.GenerateUUID()),
		UserID:         user.UserID,
		Token:          tokenString,
		ExpirationDate: expirationDate,
	}, nil
}

func (g *Generator) getExpirationDate() time.Time {
	return time.Now().Add(ExpirationDuration)
}
