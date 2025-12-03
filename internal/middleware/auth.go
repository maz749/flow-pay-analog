package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/maz749/flow-pay-analog/internal/repository"
	"github.com/maz749/flow-pay-analog/pkg/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "userRole"

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

func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(UserRoleKey).(string)
	return role, ok
}

// RequireAdmin middleware ensures that only users with admin role can access the route
func RequireAdmin(userRepo *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, "User not found")
				return
			}

			if user.Role != "admin" {
				utils.SendError(w, http.StatusForbidden, "Admin access required")
				return
			}

			// Add user role to context
			ctx := context.WithValue(r.Context(), UserRoleKey, user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole middleware ensures that users have one of the specified roles
func RequireRole(userRepo *repository.UserRepository, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r)
			if !ok {
				utils.SendError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			user, err := userRepo.GetByID(userID)
			if err != nil {
				utils.SendError(w, http.StatusUnauthorized, "User not found")
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if user.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				utils.SendError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			// Add user role to context
			ctx := context.WithValue(r.Context(), UserRoleKey, user.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
