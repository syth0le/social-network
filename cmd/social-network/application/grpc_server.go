package application

import (
	xservers "github.com/syth0le/gopnik/servers"

	"github.com/syth0le/social-network/internal/handler/internalapi"
	inpb "github.com/syth0le/social-network/proto/internalapi"
)

func (a *App) newInternalGRPCServer(env *env) *xservers.GRPCServer {
	server := xservers.NewGRPCServer(
		a.Config.InternalGRPCServer,
		a.Logger,
		xservers.GRPCWithServerName("internal grpc api"),
	)

	inpb.RegisterAuthServiceServer(server.Server, &internalapi.AuthHandler{AuthService: env.authenticationService})

	return server
}
