package token

import (
	"fmt"
	"social-network/internal/model"
)

type Generator struct {
	SaltValue string
}

func (g *Generator) GenerateToken(user *model.User) (string, error) {
	if user == nil {
		return "", fmt.Errorf("user cannot empty")
	}
	return "", nil
}
