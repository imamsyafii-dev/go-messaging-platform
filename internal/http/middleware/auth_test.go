package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/middleware"
)

func TestRequireAuthRejectsMissingToken(t *testing.T) {
	tokenService := auth.NewTokenService("test-secret", time.Hour)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})

	handler := middleware.RequireAuth(tokenService)(next)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestRequireAuthAddsUserID(t *testing.T) {
	tokenService := auth.NewTokenService("test-secret", time.Hour)

	token, err := tokenService.Generate("user-123")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := middleware.UserID(r.Context())
		if !ok {
			t.Fatal("user ID should be present in context")
		}

		if userID != "user-123" {
			t.Fatalf("expected user ID %q, got %q", "user-123", userID)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	handler := middleware.RequireAuth(tokenService)(next)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
}
