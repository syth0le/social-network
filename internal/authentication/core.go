package authentication

import (
	"net/http"

	"social-network/internal/service/user"
)

type Service struct {
	UserService user.Service
}

func (s Service) AuthenticationInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//s.UserService.GetUserByID(r.Body)
		// todo: check user is authenticated????? https://github.com/Grey2k/otus-ha/blob/main/backend/server/v1/api.go#L73
		next.ServeHTTP(w, r)
	})
}
