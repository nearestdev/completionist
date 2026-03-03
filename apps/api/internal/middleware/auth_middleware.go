package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/repository"
)

type ctxKey string

const ctxUserID ctxKey = "userID"
const ctxAuthToken ctxKey = "authToken"

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

			if userRepo != nil {
				if _, err := userRepo.FindByID(claims.UserID); err != nil {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
			}
			ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
			ctx = context.WithValue(ctx, ctxAuthToken, tokenStr)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}