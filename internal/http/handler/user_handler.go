package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/middleware"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/user"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService interface {
	CreateUser(
		ctx context.Context,
		email string,
		password string,
		username string,
		displayName string,
	) (sqlc.User, error)

	GetByID(
		ctx context.Context,
		id pgtype.UUID,
	) (sqlc.User, error)
}
type UserProfileService interface {
	GetByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error)
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

	createdUser, err := h.service.CreateUser(
		r.Context(),
		request.Email,
		request.Password,
		request.Username,
		request.DisplayName,
	)
	if err != nil {
		if errors.Is(err, user.ErrUserAlreadyExists) {
			http.Error(w, "email or username already exists", http.StatusConflict)
			return
		}

		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	response := userResponse{
		ID:          createdUser.ID.String(),
		Email:       createdUser.Email,
		Username:    createdUser.Username,
		DisplayName: createdUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "user identity not found", http.StatusUnauthorized)
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user identity", http.StatusUnauthorized)
		return
	}

	currentUser, err := h.service.GetByID(
		r.Context(),
		pgtype.UUID{
			Bytes: parsedID,
			Valid: true,
		},
	)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	response := userResponse{
		ID:          currentUser.ID.String(),
		Email:       currentUser.Email,
		Username:    currentUser.Username,
		DisplayName: currentUser.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
