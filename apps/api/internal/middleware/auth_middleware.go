package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
)

type ctxKey string

const ctxUserID ctxKey = "userID"
const ctxAuthToken ctxKey = "authToken"
const ctxUserRole ctxKey = "userRole"

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(ctxUserID)
	id, ok := v.(int64)
	return id, ok
}

func AuthTokenFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(ctxAuthToken)
	token, ok := v.(string)
	return token, ok
}

func UserRoleFromContext(ctx context.Context) (models.UserRole, bool) {
	v := ctx.Value(ctxUserRole)
	role, ok := v.(models.UserRole)
	return role, ok
}

func AuthMiddleware(userRepo *repository.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := ""
			if token := r.URL.Query().Get("auth_token"); token != "" {
				tokenStr = token
			} else {
				authz := r.Header.Get("Authorization")
				if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
					http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
					return
				}
				tokenStr = strings.TrimPrefix(authz, "Bearer ")
			}

			if tokenStr == "" {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			var role models.UserRole
			if userRepo != nil {
				user, err := userRepo.FindByID(claims.UserID)
				if err != nil {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				role = user.Role
			} else {
				role = models.RoleUser
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxAuthToken, tokenStr)
			ctx = context.WithValue(ctx, ctxUserRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(allowed ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := UserRoleFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, candidate := range allowed {
				if role == candidate {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
