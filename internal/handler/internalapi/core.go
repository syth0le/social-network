package internalapi

import (
	"context"
	"fmt"
	"net/http"

	xerrors "github.com/syth0le/gopnik/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/syth0le/social-network/internal/authentication"
	internalapi2 "github.com/syth0le/social-network/pkg/proto/internalapi"
)

type AuthHandler struct {
	internalapi2.UnimplementedAuthServiceServer

	AuthService authentication.Service
}

func (h *AuthHandler) ValidateToken(ctx context.Context, request *internalapi2.ValidateTokenRequest) (*emptypb.Empty, error) {
	err := h.AuthService.ValidateToken(request.Token)
	if err != nil {
		return nil, GRPCError(fmt.Errorf("validate token: %w", err))
	}

	return &emptypb.Empty{}, nil
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
