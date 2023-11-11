package utils

import "github.com/google/uuid"

const serviceNamePrefix = "snw"

func GenerateUUID() string {
	return serviceNamePrefix + uuid.New().String()
}
