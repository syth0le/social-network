package utils

import (
	"strings"

	"github.com/google/uuid"
)

const (
	serviceNamePrefix = "snw"
	userEntityPrefix  = "u"
	postEntityPrefix  = "p"
)

func GenerateUUID() string {
	return generateUID(userEntityPrefix)
}

func generateUID(entityPrefix string) string {
	return serviceNamePrefix + entityPrefix + strings.Replace(uuid.New().String(), "-", "", -1)
}

func GeneratePUID() string {
	return generateUID(postEntityPrefix)
}
