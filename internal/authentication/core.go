package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	xerrors "github.com/syth0le/gopnik/errors"
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
			s.writeError(w, fmt.Errorf("validate token: %w", err))
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

func (s *Service) writeError(w http.ResponseWriter, err error) {
	s.Logger.Sugar().Warnf("http response error: %v", err)

	w.Header().Set(headers.ContentType, "application/json")
	errorResult, ok := xerrors.FromError(err)
	if !ok {
		s.Logger.Sugar().Errorf("cannot write log message: %v", err)
		return
	}
	w.WriteHeader(errorResult.StatusCode)
	err = json.NewEncoder(w).Encode(
		map[string]any{
			"message": errorResult.Msg,
			"code":    errorResult.StatusCode,
		})

	if err != nil {
		http.Error(w, xerrors.InternalErrorMessage, http.StatusInternalServerError) // TODO: make error mapping
	}
}
