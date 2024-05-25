package utils

import (
	"strings"

	"github.com/google/uuid"
)

const (
	serviceNamePrefix  = "snw"
	userEntityPrefix   = "u"
	tokenEntityPrefix  = "t"
	postEntityPrefix   = "p"
	friendEntityPrefix = "f"
)

func GenerateUUID() string {
	return generateUID(userEntityPrefix)
}

func GenerateTUID() string {
	return generateUID(tokenEntityPrefix)
}

func generateUID(entityPrefix string) string {
	return serviceNamePrefix + entityPrefix + strings.Replace(uuid.New().String(), "-", "", -1)
}

func GeneratePUID() string {
	return generateUID(postEntityPrefix)
}

func GenerateFUID() string {
	return generateUID(friendEntityPrefix)
}
