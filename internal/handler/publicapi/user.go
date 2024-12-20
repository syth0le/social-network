package publicapi

import (
	"fmt"
	"net/http"
	"time"

	xerrors "github.com/syth0le/gopnik/errors"

	"github.com/go-chi/chi/v5"

	"github.com/syth0le/social-network/internal/model"
	"github.com/syth0le/social-network/internal/service/user"
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

// todo pagination
type userListResponse struct {
	Users []*userResponse `json:"users"`
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
	handleRequest := func() (*userListResponse, error) {
		ctx := r.Context()
		firstName := r.URL.Query().Get("first_name")
		secondName := r.URL.Query().Get("second_name")
		if firstName == "" {
			return nil, xerrors.WrapValidationError(fmt.Errorf("incorrect query args. firstName cannot by empty"))
		}
		if secondName == "" {
			return nil, xerrors.WrapValidationError(fmt.Errorf("incorrect query args. secondName cannot by empty"))
		}

		userModel, err := h.UserService.SearchUser(ctx, &user.SearchUserParams{
			FirstName:  firstName,
			SecondName: secondName,
		})
		if err != nil {
			return nil, fmt.Errorf("search user: %w", err)
		}

		return userModelsToResponse(userModel), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("search user: %w", err))
		return
	}
	writeResponse(w, response)
}

func (h *Handler) SearchTarantoolUser(w http.ResponseWriter, r *http.Request) {
	handleRequest := func() (*userListResponse, error) {
		ctx := r.Context()
		firstName := r.URL.Query().Get("first_name")
		secondName := r.URL.Query().Get("second_name")
		if firstName == "" {
			return nil, xerrors.WrapValidationError(fmt.Errorf("incorrect query args. firstName cannot by empty"))
		}
		if secondName == "" {
			return nil, xerrors.WrapValidationError(fmt.Errorf("incorrect query args. secondName cannot by empty"))
		}

		userModel, err := h.UserService.SearchTarantoolUser(ctx, &user.SearchTarantoolUserParams{
			FirstName:  firstName,
			SecondName: secondName,
		})
		if err != nil {
			return nil, fmt.Errorf("search tarantool user: %w", err)
		}

		return tarantoolUserModelsToResponse(userModel), nil
	}

	response, err := handleRequest()
	if err != nil {
		h.writeError(r.Context(), w, fmt.Errorf("search tarantool user: %w", err))
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

func userModelsToResponse(userModels []*model.User) *userListResponse {
	users := make([]*userResponse, 0)
	for _, userModel := range userModels {
		users = append(users, userModelToResponse(userModel))
	}
	return &userListResponse{
		Users: users,
	}
}

func tarantoolUserModelToResponse(user model.TarantoolUser) *userResponse {
	return &userResponse{
		UserID:     user.UserID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Sex:        user.Sex,
		Biography:  user.Biography,
		City:       user.City,
	}
}

func tarantoolUserModelsToResponse(userModels []model.TarantoolUser) *userListResponse {
	users := make([]*userResponse, 0)
	for _, userModel := range userModels {
		users = append(users, tarantoolUserModelToResponse(userModel))
	}
	return &userListResponse{
		Users: users,
	}
}
