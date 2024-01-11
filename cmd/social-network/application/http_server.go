package application

import (
	"social-network/internal/handler/publicapi"

	"github.com/go-chi/chi/v5"
)

func (a *App) newHTTPServer(env *env) *HTTPServerWrapper {
	return NewHTTPServerWrapper(
		a.Logger,
		WithAdminServer(a.Config.AdminServer),
		WithPublicServer(a.Config.PublicServer, a.publicMux(env)),
	)
}

func (a *App) publicMux(env *env) *chi.Mux {
	mux := chi.NewMux()

	handler := publicapi.Handler{
		Logger:      a.Logger,
		UserService: env.userService,
	}

	mux.Post("/login", handler.Login)
	mux.Post("/user/register", handler.Register)
	mux.Route("/user", func(r chi.Router) {
		r.Use(env.authenticationService.AuthenticationInterceptor)

		r.Get("/{userID}", handler.GetUserByID)

		r.Get("/search", handler.SearchUser)
	})

	return mux
}
