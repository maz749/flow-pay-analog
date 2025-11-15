package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.SendError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Expected format: Bearer <token>
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.SendError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			token := parts[1]
			claims, err := utils.ValidateToken(token, jwtSecret)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(UserIDKey).(int)
	return userID, ok
}
