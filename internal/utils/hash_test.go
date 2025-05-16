package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckPasswordHashSuccess(t *testing.T) {
	var tableData = map[string]struct{ password, expected string }{
		"success test data 1": {"324234", "$2a$14$2DpOdTKDyf3TXhUPmWyFxeS50C3WxlAuOIkoUHYSpF4F12M/kEU.i"},
		"success test data 2": {"222222", "$2a$14$2DpOdTKDyfjhghjPmWyFxeS50C3WxlAuOIkoUHYSpF4F12M/kEU.i"},
	}

	for name, data := range tableData {
		t.Run(name, func(t *testing.T) {
			result, err := HashPassword(data.password)

			require.NoError(t, err)
			require.Equal(t, len(result), len(data.expected))
		})
	}
}

func TestCheckPasswordHashFailed(t *testing.T) {
	var tableData = map[string]string{
		"failed test data 1": "",
	}

	for name, data := range tableData {
		t.Run(name, func(t *testing.T) {
			_, err := HashPassword(data)

			require.Error(t, err)
		})
	}
}
