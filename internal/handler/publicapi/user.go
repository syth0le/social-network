package publicapi

import (
	"fmt"
	"net/http"
	"social-network/internal/model"
	"social-network/internal/service/user"
	"social-network/internal/utils"
	"time"

	"github.com/go-chi/chi/v5"
)

type jwtTokenResponse struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

type userResponse struct {
	UserID     string `json:"user_id"`
	Username   string `json:"username"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Sex        string `json:"sex"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*jwtTokenResponse, error) {
		ctx := r.Context()

		request, err := parseJSONRequest[loginRequest](r)
		if err != nil {
			return nil, fmt.Errorf("parse json request: %w", err)
		}

		tokenModel, err := h.UserService.Login(ctx, &user.LoginParams{
			Username: request.Username,
			Password: request.Password,
		})
		if err != nil {
			return nil, fmt.Errorf("login: %w", err)
		}

		return &jwtTokenResponse{
			Token:  tokenModel.Token,
			UserID: tokenModel.UserID.String(),
		}, nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("login: %w", err))
		return
	}
	writeResponse(w, response)
}

type registerRequest struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	FirstName  string    `json:"first_name"`
	SecondName string    `json:"second_name"`
	Age        int       `json:"age"`
	Sex        string    `json:"sex"`
	Birthdate  time.Time `json:"birthdate"`
	Biography  string    `json:"biography"`
	City       string    `json:"city"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*jwtTokenResponse, error) {
		ctx := r.Context()

		request, err := parseJSONRequest[registerRequest](r)
		if err != nil {
			return nil, fmt.Errorf("parse json request: %w", err)
		}

		tokenModel, err := h.UserService.Register(ctx, &user.RegisterParams{
			Username:   request.Username,
			Password:   request.Password,
			FirstName:  request.FirstName,
			SecondName: request.SecondName,
			Age:        request.Age,
			Sex:        request.Sex,
			Birthdate:  request.Birthdate,
			Biography:  request.Biography,
			City:       request.City,
		})
		if err != nil {
			return nil, fmt.Errorf("register: %w", err)
		}

		return &jwtTokenResponse{
			Token:  tokenModel.Token,
			UserID: tokenModel.UserID.String(),
		}, nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("register: %w", err))
		return
	}
	writeResponse(w, response)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*userResponse, error) {
		ctx := r.Context()
		userID := chi.URLParamFromCtx(ctx, "userID")

		userModel, err := h.UserService.GetUserByID(ctx, &user.GetUserByIDParams{
			UserID: model.UserID(userID),
		})
		if err != nil {
			return nil, fmt.Errorf("get user by id: %w", err)
		}

		return userModelToResponse(userModel), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("get user by id: %w", err))
		return
	}
	writeResponse(w, response)
}

func (h *Handler) SearchUser(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*userResponse, error) {
		ctx := r.Context()
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		if firstName == "" {
			return nil, utils.WrapValidationError(fmt.Errorf("incorrect query args. first_name cannot by empty"))
		}
		if lastName == "" {
			return nil, utils.WrapValidationError(fmt.Errorf("incorrect query args. first_name cannot by empty"))
		}

		userModel, err := h.UserService.SearchUser(ctx, &user.SearchUserParams{
			FirstName: firstName,
			LastName:  lastName,
		})
		if err != nil {
			return nil, fmt.Errorf("search user: %w", err)
		}

		return userModelToResponse(userModel), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("search user: %w", err))
		return
	}
	writeResponse(w, response)
}

func userModelToResponse(user *model.User) *userResponse {
	return &userResponse{
		UserID:     user.UserID.String(),
		Username:   user.Username,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Sex:        user.Sex,
		Birthdate:  user.Birthdate.String(),
		Biography:  user.Biography,
		City:       user.City,
	}
}
