package auth

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/syth0le/dialog-service/internal/model"
)

type ClientMock struct {
	logger *zap.Logger
}

func NewClientMock(logger *zap.Logger) *ClientMock {
	return &ClientMock{
		logger: logger,
	}
}

func (m *ClientMock) AuthenticationInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Debug("authenticated through mock service")
		ctx := context.WithValue(r.Context(), UserIDValue, model.UserID("mock_user"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
