package application

import (
	xservers "github.com/syth0le/gopnik/servers"

	"social-network/internal/handler/publicapi"

	"github.com/go-chi/chi/v5"
)

func (a *App) newHTTPServer(env *env) *xservers.HTTPServerWrapper {
	return xservers.NewHTTPServerWrapper(
		a.Logger,
		xservers.WithAdminServer(a.Config.AdminServer),
		xservers.WithPublicServer(a.Config.PublicServer, a.publicMux(env)),
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
