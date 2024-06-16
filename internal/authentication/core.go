package authentication

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/syth0le/social-network/internal/service/user"
	"github.com/syth0le/social-network/internal/token"
)

const authHeader = "Authorization"
const UserIDValue = "userID"

type Service struct {
	Logger       *zap.Logger
	UserService  user.Service
	TokenManager *token.Manager
}

func (s *Service) AuthenticationInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(authHeader)

		userID, err := s.TokenManager.ValidateToken(authToken)
		if err != nil {
			s.Logger.Sugar().Errorf("validate token: %v", err)
			return
		}

		// todo: check token valid in db? (мы могли сделать инвалидейт)
		// user, err := s.UserService.GetUserByTokenAndID(r.Context(), &user.GetUserByTokenAndIDParams{Token: authToken})

		ctx := context.WithValue(r.Context(), UserIDValue, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) ValidateToken(authToken string) error {
	_, err := s.TokenManager.ValidateToken(authToken)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}

	return nil
}
