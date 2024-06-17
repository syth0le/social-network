package internalapi

import (
	"context"
	"fmt"
	"net/http"

	xerrors "github.com/syth0le/gopnik/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/syth0le/social-network/internal/authentication"
	"github.com/syth0le/social-network/proto/internalapi"
)

type AuthHandler struct {
	internalapi.UnimplementedAuthServiceServer

	AuthService authentication.Service
}

// todo: validate token raise 400 instead of 403
func (h *AuthHandler) ValidateToken(ctx context.Context, request *internalapi.ValidateTokenRequest) (*internalapi.ValidateTokenResponse, error) {
	userID, err := h.AuthService.ValidateToken(request.Token)
	if err != nil {
		return nil, GRPCError(fmt.Errorf("validate token: %w", err))
	}

	return &internalapi.ValidateTokenResponse{UserId: userID.String()}, nil
}

// GRPCError todo: move to gopnik and create server interceptor
func GRPCError(err error) error {
	resError, ok := xerrors.FromError(err)
	if !ok {
		return err
	}

	switch resError.StatusCode {
	case http.StatusForbidden:
		return status.Error(codes.PermissionDenied, resError.Msg)
	default:
		return status.Errorf(codes.Internal, resError.Msg)
	}
}
