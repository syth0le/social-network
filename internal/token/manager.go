package token

import (
	"fmt"
	"time"

	xerrors "github.com/syth0le/gopnik/errors"

	"github.com/syth0le/social-network/cmd/social-network/configuration"
	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

const ExpirationDuration = time.Hour * 24 // todo: get expiration time from settings

type Manager struct {
	saltValue  []byte
	serverName string
}

func NewManager(config configuration.ApplicationConfig) *Manager {
	return &Manager{
		saltValue:  []byte(config.SaltValue),
		serverName: config.App,
	}
}

func (m *Manager) GenerateToken(user *model.User) (*model.Token, error) {
	if user == nil {
		return nil, fmt.Errorf("user cannot be empty")
	}

	expirationDate := m.getExpirationDate()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": m.serverName,
		"uid": user.UserID,
		"usr": user.Username,
		"exp": expirationDate.Unix(),
	})

	tokenString, err := token.SignedString(m.saltValue)
	if err != nil {
		return nil, fmt.Errorf("signed string: %w", err)
	}

	return &model.Token{
		TokenID:        model.TokenID(utils.GenerateTUID()),
		UserID:         user.UserID,
		Token:          tokenString,
		ExpirationDate: expirationDate,
	}, nil
}

func (m *Manager) ValidateToken(tokenString string) (model.UserID, error) {
	claimsMap, err := m.parseToken(tokenString)
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	id, ok := claimsMap["uid"]
	if !ok {
		return "", xerrors.WrapForbiddenError(fmt.Errorf("not found uid"), "parse token error")
	}

	return model.UserID(id.(string)), nil
}

func (m *Manager) parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return m.saltValue, nil
	})
	if err != nil {
		return nil, xerrors.WrapForbiddenError(fmt.Errorf("parsing token err: %w", err), "parse token error")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, xerrors.WrapForbiddenError(fmt.Errorf("claims token err or token is not valid"), "parse token error")
}

func (m *Manager) getExpirationDate() time.Time {
	return time.Now().Add(ExpirationDuration)
}
