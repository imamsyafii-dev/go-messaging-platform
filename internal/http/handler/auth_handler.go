package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/user"
)

type AuthUserService interface {
	Authenticate(
		ctx context.Context,
		email string,
		password string,
	) (sqlc.User, error)
}

type AuthHandler struct {
	userService  AuthUserService
	tokenService *auth.TokenService
}

func NewAuthHandler(
	userService AuthUserService,
	tokenService *auth.TokenService,
) *AuthHandler {
	return &AuthHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request loginRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	currentUser, err := h.userService.Authenticate(
		r.Context(),
		request.Email,
		request.Password,
	)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		http.Error(w, "failed to authenticate", http.StatusInternalServerError)
		return
	}

	accessToken, err := h.tokenService.Generate(currentUser.ID.String())
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(loginResponse{
		AccessToken: accessToken,
	})
}
