package publicapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/internal/service/user"

	"github.com/go-http-utils/headers"
)

type Handler struct {
	Logger      log.Logger
	UserService user.Service
}

func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, err error) {
	// TODO
}

func writeResponse(w http.ResponseWriter, response any) {
	w.Header().Set(headers.ContentType, "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError) // TODO: make error mapping
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
