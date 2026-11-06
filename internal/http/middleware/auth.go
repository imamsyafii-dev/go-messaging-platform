package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
)

type contextKey string

const userIDKey contextKey = "user_id"

func RequireAuth(tokenService *auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			parts := strings.SplitN(header, " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			userID, err := tokenService.Parse(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
