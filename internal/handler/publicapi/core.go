package publicapi

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"social-network/internal/service/user"
	"social-network/internal/utils"

	"github.com/go-http-utils/headers"
)

type Handler struct {
	Logger      *zap.Logger
	UserService user.Service
}

func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error) {
	h.Logger.Sugar().Warnf("http response error: %v", err)

	w.Header().Set(headers.ContentType, "application/json")
	errorResult, ok := utils.FromError(err)
	if !ok {
		h.Logger.Error("cannot write log message")
		return
	}
	err = json.NewEncoder(w).Encode(
		map[string]any{
			"message": errorResult.Msg,
			"code":    errorResult.StatusCode,
		})
	w.WriteHeader(errorResult.StatusCode)
	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError) // TODO: make error mapping
	}
}

func writeResponse(w http.ResponseWriter, response any) {
	w.Header().Set(headers.ContentType, "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, utils.InternalErrorMessage, http.StatusInternalServerError) // TODO: make error mapping
	}
}

func parseJSONRequest[T loginRequest | registerRequest](r *http.Request) (*T, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = fmt.Errorf("read body: %w", err)
		return nil, err // TODO: wrap error to ValidationError
	}

	var request T
	err = json.Unmarshal(body, &request)
	if err != nil {
		err = fmt.Errorf("unmarshal request body: %w", err)
		return nil, err // TODO: wrap error to ValidationError
	}
	return &request, nil
}
