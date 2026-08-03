package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/security"
)

func TestAuthMiddleware(t *testing.T) {
	manager := security.NewTokenManager(
		"abcdefghijklmnopqrstuvwxyz-1234567890",
		time.Hour,
	)

	token, err := manager.Generate(&models.User{
		ID:       1,
		FullName: "Administrator",
		Email:    "admin@test.local",
		Role:     models.RoleAdmin,
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	handler := Auth(
		manager,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := AuthUserFromContext(r.Context())
			if !ok || user.ID != 1 {
				t.Fatal("auth user tidak tersedia")
			}
			w.WriteHeader(http.StatusNoContent)
		}),
	)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d", recorder.Code)
	}
}

func TestRequireRolesRejectsWrongRole(t *testing.T) {
	manager := security.NewTokenManager(
		"abcdefghijklmnopqrstuvwxyz-1234567890",
		time.Hour,
	)

	token, err := manager.Generate(&models.User{
		ID:       2,
		FullName: "Karyawan",
		Email:    "employee@test.local",
		Role:     models.RoleEmployee,
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	handler := Auth(
		manager,
		RequireRoles(models.RoleAdmin)(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}),
		),
	)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d", recorder.Code)
	}
}
