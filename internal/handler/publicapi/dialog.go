package publicapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syth0le/social-network/internal/authentication"
	"github.com/syth0le/social-network/internal/model"
)

type createDialogRequest struct {
	ParticipantsIDs []string `json:"participants_ids"`
}

type dialogResponse struct {
	DialogID        string   `json:"dialog_id"`
	ParticipantsIDs []string `json:"participants_ids"`
}

// deprecatedAPI

func (h *Handler) CreateDialog(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*model.Dialog, error) {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return nil, fmt.Errorf("cannot recognize userID")
		}

		request, err := parseJSONRequest[createDialogRequest](r)
		if err != nil {
			return nil, fmt.Errorf("parse json request: %w", err)
		}

		participantsIDs := make([]model.UserID, len(request.ParticipantsIDs))
		for idx, item := range request.ParticipantsIDs {
			participantsIDs[idx] = model.UserID(item)
		}

		dialogModel, err := h.DialogClient.CreateDialog(ctx, userIDStr.(model.UserID), participantsIDs)
		if err != nil {
			return nil, fmt.Errorf("create dialog: %w", err)
		}

		return dialogModel, nil
	}

	resp, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("create dialog: %w", err))
		return
	}

	writeResponse(w, dialogModelToResponse(resp))
}

type createMessageRequest struct {
	DialogID string `json:"dialog_id"`
	Text     string `json:"text"`
}

func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return fmt.Errorf("cannot recognize userID")
		}

		request, err := parseJSONRequest[createMessageRequest](r)
		if err != nil {
			return fmt.Errorf("parse json request: %w", err)
		}

		err = h.DialogClient.CreateMessage(ctx, model.DialogID(request.DialogID), userIDStr.(model.UserID), request.Text)
		if err != nil {
			return fmt.Errorf("create message: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("create message: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

type messageResponse struct {
	ID       string `json:"id"`
	DialogID string `json:"dialog_id"`
	Text     string `json:"text"`
	SenderID string `json:"sender_id"`
}

type messageListResponse struct {
	Messages []*messageResponse `json:"messages"`
}

func (h *Handler) GetDialogMessages(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*messageListResponse, error) {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return nil, fmt.Errorf("cannot recognize userID")
		}

		dialogID := chi.URLParamFromCtx(ctx, "dialogID")

		messages, err := h.DialogClient.GetDialogMessages(ctx, model.DialogID(dialogID), userIDStr.(model.UserID))
		if err != nil {
			return nil, fmt.Errorf("get dialog: %w", err)
		}

		return messageModelsToResponse(messages), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("get dialog: %w", err))
		return
	}

	writeResponse(w, response)
}

func dialogModelToResponse(dialog *model.Dialog) *dialogResponse {
	participants := make([]string, len(dialog.ParticipantsIDs))
	for idx, item := range dialog.ParticipantsIDs {
		participants[idx] = item.String()
	}

	return &dialogResponse{
		DialogID:        dialog.ID.String(),
		ParticipantsIDs: participants,
	}
}

func messageModelToResponse(message *model.Message) *messageResponse {
	return &messageResponse{
		ID:       message.ID.String(),
		DialogID: message.DialogID.String(),
		Text:     message.Text,
		SenderID: message.SenderID.String(),
	}
}

func messageModelsToResponse(messageModels []*model.Message) *messageListResponse {
	messages := make([]*messageResponse, 0)
	for _, messageModel := range messageModels {
		messages = append(messages, messageModelToResponse(messageModel))
	}

	return &messageListResponse{
		Messages: messages,
	}
}
