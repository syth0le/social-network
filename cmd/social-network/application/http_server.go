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
		PostService: env.postService,
	}

	mux.Post("/login", handler.Login)
	mux.Post("/user/register", handler.Register)
	mux.Route("/user", func(r chi.Router) {
		r.Get("/{userID}", handler.GetUserByID)

		r.Get("/search", handler.SearchUser)
	})

	mux.Route("/friend", func(r chi.Router) {
		r.Use(env.authenticationService.AuthenticationInterceptor)

		r.Get("/{userID}", handler.ListFriends)
		r.Put("/set/{userID}", handler.SetFriendRequest)
		r.Put("/delete/{userID}", handler.DeleteFriend)
	})

	mux.Route("/post", func(r chi.Router) {
		r.Use(env.authenticationService.AuthenticationInterceptor)

		r.Post("/", handler.CreatePost)
		r.Get("/{postID}", handler.GetPostByID)
		r.Patch("/{postID}", handler.UpdatePost)
		r.Delete("/{postID}", handler.DeletePost)

		r.Get("/feed", handler.GetFeed)
	})

	return mux
}
