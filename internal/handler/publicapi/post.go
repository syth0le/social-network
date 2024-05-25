package publicapi

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"social-network/internal/authentication"
	"social-network/internal/model"
	"social-network/internal/service/post"
)

type createPostRequest struct {
	Text string `json:"text"`
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return fmt.Errorf("cannot recognize userID")
		}

		request, err := parseJSONRequest[createPostRequest](r)
		if err != nil {
			return fmt.Errorf("parse json request: %w", err)
		}

		err = h.PostService.Create(ctx, &post.CreatePostParams{
			UserID: userIDStr.(model.UserID),
			Text:   request.Text,
		})
		if err != nil {
			return fmt.Errorf("post create: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("create post: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

type postResponse struct {
	ID       string `json:"id"`
	Text     string `json:"text"`
	Author   string `json:"author"`
	AuthorID string `json:"author_id"`
}

func (h *Handler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*postResponse, error) {
		ctx := r.Context()
		postID := chi.URLParamFromCtx(ctx, "postID")

		postModel, err := h.PostService.GetPostByID(ctx, &post.GetPostByIDParams{
			PostID: model.PostID(postID),
		})
		if err != nil {
			return nil, fmt.Errorf("get post by id: %w", err)
		}

		return postModelToResponse(postModel), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("get post by id: %w", err))
		return
	}
	writeResponse(w, response)
}

type updatePostRequest struct {
	Text string `json:"text"`
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()
		postID := chi.URLParamFromCtx(ctx, "postID")

		request, err := parseJSONRequest[updatePostRequest](r)
		if err != nil {
			return fmt.Errorf("parse json request: %w", err)
		}

		err = h.PostService.Update(ctx, &post.UpdatePostParams{
			PostID: model.PostID(postID),
			Text:   request.Text,
		})
		if err != nil {
			return fmt.Errorf("update post: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("update post: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() error {
		ctx := r.Context()
		postID := chi.URLParamFromCtx(ctx, "postID")

		err := h.PostService.Delete(ctx, &post.DeletePostParams{
			PostID: model.PostID(postID),
		})
		if err != nil {
			return fmt.Errorf("delete post: %w", err)
		}

		return nil
	}

	err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("delete post: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}

type getFeedResponse struct {
	Posts []*postResponse `json:"posts"`
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*getFeedResponse, error) {
		ctx := r.Context()

		userIDStr := ctx.Value(authentication.UserIDValue)
		if userIDStr == "" {
			return nil, fmt.Errorf("cannot recognize userID")
		}

		posts, err := h.PostService.GetFeed(ctx, &post.GetFeedParams{
			FollowerID: userIDStr.(model.UserID),
		})
		if err != nil {
			return nil, fmt.Errorf("get feed: %w", err)
		}

		return postModelsToResponse(posts), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("get feed: %w", err))
		return
	}
	writeResponse(w, response)
}

func postModelToResponse(post *model.Post) *postResponse {
	return &postResponse{
		ID:       post.ID.String(),
		Text:     post.Text,
		Author:   post.Author,
		AuthorID: post.AuthorID.String(),
	}
}

func postModelsToResponse(postModels []*model.Post) *getFeedResponse {
	posts := make([]*postResponse, 0)
	for _, postModel := range postModels {
		posts = append(posts, postModelToResponse(postModel))
	}
	return &getFeedResponse{
		Posts: posts,
	}
}
