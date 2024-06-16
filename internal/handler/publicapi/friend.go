package publicapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/syth0le/social-network/internal/authentication"
	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/service/friend"
)

type friendResponse struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
}

type friendListResponse struct {
	Users []*friendResponse `json:"users"`
}

func (h *Handler) ListFriends(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*friendListResponse, error) {
		ctx := r.Context()

		userID := chi.URLParamFromCtx(ctx, "userID")

		friends, err := h.FriendService.ListFriends(ctx, &friend.ListFriendsParams{
			UserID: model.UserID(userID),
		})
		if err != nil {
			return nil, fmt.Errorf("list friends: %w", err)
		}

		return friendModelsToResponse(friends), nil
	}

	resp, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("set friend: %w", err))
		return
	}

	writeResponse(w, resp)
}

func (h *Handler) SetFriendRequest(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return fmt.Errorf("cannot recognize userID")
		}

		authorID := chi.URLParamFromCtx(ctx, "userID")

		err := h.FriendService.AddFriend(ctx, &friend.AddFriendParams{
			AuthorID:   userIDStr.(model.UserID),
			FollowerID: model.UserID(authorID),
		})
		if err != nil {
			return fmt.Errorf("set friend request: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("set friend request: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return fmt.Errorf("cannot recognize userID")
		}

		authorID := chi.URLParamFromCtx(ctx, "userID")

		err := h.FriendService.DeleteFriend(ctx, &friend.DeleteFriendParams{
			AuthorID:   userIDStr.(model.UserID),
			FollowerID: model.UserID(authorID),
		})
		if err != nil {
			return fmt.Errorf("set friend: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("delete friend: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func friendModelToResponse(user *model.Friend) *friendResponse {
	return &friendResponse{
		UserID:     user.UserID.String(),
		Username:   user.Username,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
	}
}

func friendModelsToResponse(friendModels []*model.Friend) *friendListResponse {
	users := make([]*friendResponse, 0)
	for _, userModel := range friendModels {
		users = append(users, friendModelToResponse(userModel))
	}

	return &friendListResponse{
		Users: users,
	}
}
