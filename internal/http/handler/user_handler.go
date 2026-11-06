package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		email string,
		password string,
		username string,
		displayName string,
	) (sqlc.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type createUserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type userResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	request.Email = strings.TrimSpace(request.Email)
	request.Username = strings.TrimSpace(request.Username)
	request.DisplayName = strings.TrimSpace(request.DisplayName)

	if request.Email == "" ||
		request.Password == "" ||
		request.Username == "" ||
		request.DisplayName == "" {
		http.Error(w, "email, password, username, and display_name are required", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(
		r.Context(),
		request.Email,
		request.Password,
		request.Username,
		request.DisplayName,
	)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	response := userResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
