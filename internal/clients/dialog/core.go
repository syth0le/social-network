package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-http-utils/headers"
	xerrors "github.com/syth0le/gopnik/errors"
	inpb "github.com/syth0le/social-network/proto/internalapi"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/syth0le/dialog-service/internal/model"
)

const authHeader = "Authorization"
const UserIDValue = "userID"

type Client interface {
	AuthenticationInterceptor(next http.Handler) http.Handler
}

type ClientImpl struct {
	logger *zap.Logger
	client inpb.AuthServiceClient
}

func NewAuthImpl(logger *zap.Logger, conn *grpc.ClientConn) *ClientImpl {
	return &ClientImpl{
		logger: logger,
		client: inpb.NewAuthServiceClient(conn),
	}
}

func (c *ClientImpl) AuthenticationInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(authHeader)

		resp, err := c.client.ValidateToken(r.Context(), &inpb.ValidateTokenRequest{Token: authToken})
		if err != nil {
			c.writeError(w, xerrors.WrapForbiddenError(fmt.Errorf("validate token: %w", err), "token validation failed"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDValue, model.UserID(resp.UserId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *ClientImpl) writeError(w http.ResponseWriter, err error) {
	c.logger.Sugar().Warnf("http response error: %v", err)

	w.Header().Set(headers.ContentType, "application/json")
	errorResult, ok := xerrors.FromError(err)
	if !ok {
		c.logger.Sugar().Errorf("cannot write log message: %v", err)
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

func (c *ClientImpl) writeGRPCError(w http.ResponseWriter, err error) {
	c.logger.Sugar().Warnf("http response error: %v", err)

	w.Header().Set(headers.ContentType, "application/json")
	errorResult, ok := xerrors.FromError(err)
	if !ok {
		c.logger.Sugar().Errorf("cannot write log message: %v", err)
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
