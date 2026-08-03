package security

import (
	"errors"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

func TestGenerateAndParseToken(t *testing.T) {
	manager := NewTokenManager(
		"abcdefghijklmnopqrstuvwxyz-1234567890",
		time.Hour,
	)

	fixedNow := time.Date(2026, 7, 17, 10, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return fixedNow }

	user := &models.User{
		ID:       10,
		FullName: "Nico",
		Email:    "nico@example.com",
		Role:     models.RoleAdmin,
	}

	token, err := manager.Generate(user)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.Subject != user.ID || claims.Role != models.RoleAdmin {
		t.Fatalf("claims tidak sesuai: %+v", claims)
	}
}

func TestExpiredToken(t *testing.T) {
	manager := NewTokenManager(
		"abcdefghijklmnopqrstuvwxyz-1234567890",
		time.Minute,
	)

	fixedNow := time.Date(2026, 7, 17, 10, 0, 0, 0, time.UTC)
	manager.now = func() time.Time { return fixedNow }

	token, err := manager.Generate(&models.User{
		ID:       1,
		FullName: "User",
		Email:    "user@example.com",
		Role:     models.RoleEmployee,
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	manager.now = func() time.Time {
		return fixedNow.Add(2 * time.Minute)
	}

	_, err = manager.Parse(token)
	if !errors.Is(err, ErrExpiredToken) {
		t.Fatalf("Parse() error = %v", err)
	}
}
