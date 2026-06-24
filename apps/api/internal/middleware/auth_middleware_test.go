package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/GATEOPENERZ/completionist-api/internal/auth"
	"github.com/GATEOPENERZ/completionist-api/internal/models"
)

func recordingHandler(ran *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*ran = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequireRoles(t *testing.T) {
	tests := []struct {
		name       string
		allowed    []models.UserRole
		setRole    bool
		role       models.UserRole
		wantStatus int
		wantNext   bool
	}{
		{
			name:       "allowed role passes through",
			allowed:    []models.UserRole{models.RoleAdmin},
			setRole:    true,
			role:       models.RoleAdmin,
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			name:       "role in allowed set among several",
			allowed:    []models.UserRole{models.RoleMember, models.RoleAdmin},
			setRole:    true,
			role:       models.RoleMember,
			wantStatus: http.StatusOK,
			wantNext:   true,
		},
		{
			name:       "disallowed role is forbidden",
			allowed:    []models.UserRole{models.RoleAdmin},
			setRole:    true,
			role:       models.RoleUser,
			wantStatus: http.StatusForbidden,
			wantNext:   false,
		},
		{
			name:       "no role in context is unauthorized",
			allowed:    []models.UserRole{models.RoleAdmin},
			setRole:    false,
			wantStatus: http.StatusUnauthorized,
			wantNext:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextRan := false
			handler := RequireRoles(tt.allowed...)(recordingHandler(&nextRan))

			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if tt.setRole {
				req = req.WithContext(context.WithValue(req.Context(), ctxUserRole, tt.role))
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if nextRan != tt.wantNext {
				t.Errorf("next ran = %v, want %v", nextRan, tt.wantNext)
			}
		})
	}
}

func TestUserIDFromContext(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxUserID, int64(42))
		got, ok := UserIDFromContext(ctx)
		if !ok {
			t.Fatal("ok = false, want true")
		}
		if got != 42 {
			t.Errorf("id = %d, want 42", got)
		}
	})

	t.Run("absent", func(t *testing.T) {
		got, ok := UserIDFromContext(context.Background())
		if ok {
			t.Errorf("ok = true, want false")
		}
		if got != 0 {
			t.Errorf("id = %d, want 0", got)
		}
	})
}

func TestAuthTokenFromContext(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxAuthToken, "abc.def.ghi")
		got, ok := AuthTokenFromContext(ctx)
		if !ok {
			t.Fatal("ok = false, want true")
		}
		if got != "abc.def.ghi" {
			t.Errorf("token = %q, want %q", got, "abc.def.ghi")
		}
	})

	t.Run("absent", func(t *testing.T) {
		got, ok := AuthTokenFromContext(context.Background())
		if ok {
			t.Errorf("ok = true, want false")
		}
		if got != "" {
			t.Errorf("token = %q, want empty", got)
		}
	})
}

func TestUserRoleFromContext(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), ctxUserRole, models.RoleAdmin)
		got, ok := UserRoleFromContext(ctx)
		if !ok {
			t.Fatal("ok = false, want true")
		}
		if got != models.RoleAdmin {
			t.Errorf("role = %q, want %q", got, models.RoleAdmin)
		}
	})

	t.Run("absent", func(t *testing.T) {
		got, ok := UserRoleFromContext(context.Background())
		if ok {
			t.Errorf("ok = true, want false")
		}
		if got != "" {
			t.Errorf("role = %q, want empty", got)
		}
	})
}

func TestAuthMiddlewareNilRepo(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")

	const userID = int64(7)
	token, err := auth.GenerateJWT(userID, time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	t.Run("valid bearer token sets user role and id", func(t *testing.T) {
		var gotID int64
		var gotIDOk bool
		var gotRole models.UserRole
		var gotRoleOk bool
		var gotToken string

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotID, gotIDOk = UserIDFromContext(r.Context())
			gotRole, gotRoleOk = UserRoleFromContext(r.Context())
			gotToken, _ = AuthTokenFromContext(r.Context())
			w.WriteHeader(http.StatusOK)
		})

		handler := AuthMiddleware(nil)(next)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rec.Code)
		}
		if !gotIDOk || gotID != userID {
			t.Errorf("user id = %d (ok=%v), want %d", gotID, gotIDOk, userID)
		}
		if !gotRoleOk || gotRole != models.RoleUser {
			t.Errorf("role = %q (ok=%v), want %q", gotRole, gotRoleOk, models.RoleUser)
		}
		if gotToken != token {
			t.Errorf("token = %q, want %q", gotToken, token)
		}
	})

	t.Run("valid token via query param", func(t *testing.T) {
		ran := false
		handler := AuthMiddleware(nil)(recordingHandler(&ran))

		req := httptest.NewRequest(http.MethodGet, "/me?auth_token="+token, nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want 200", rec.Code)
		}
		if !ran {
			t.Error("next did not run")
		}
	})

	rejectTests := []struct {
		name      string
		setHeader bool
		header    string
	}{
		{name: "missing authorization header", setHeader: false},
		{name: "non-bearer header", setHeader: true, header: "Basic " + token},
		{name: "invalid token", setHeader: true, header: "Bearer not-a-real-token"},
	}

	for _, tt := range rejectTests {
		t.Run(tt.name, func(t *testing.T) {
			ran := false
			handler := AuthMiddleware(nil)(recordingHandler(&ran))

			req := httptest.NewRequest(http.MethodGet, "/me", nil)
			if tt.setHeader {
				req.Header.Set("Authorization", tt.header)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want 401", rec.Code)
			}
			if ran {
				t.Error("next ran, want it not to run")
			}
		})
	}
}
