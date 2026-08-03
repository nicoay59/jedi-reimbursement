package services

import (
	"context"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/security"
)

type fakeUserFinder struct {
	user *models.User
	err  error
}

func (f fakeUserFinder) FindByEmail(
	context.Context,
	string,
) (*models.User, error) {
	return f.user, f.err
}

func (f fakeUserFinder) FindByID(
	context.Context,
	int64,
) (*models.User, error) {
	return f.user, f.err
}

func TestLoginSuccess(t *testing.T) {
	passwordHash, err := security.HashPassword("Password123!")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	service := NewAuthService(
		fakeUserFinder{
			user: &models.User{
				ID:           1,
				FullName:     "Administrator",
				Email:        "admin@test.local",
				PasswordHash: passwordHash,
				Role:         models.RoleAdmin,
				IsActive:     true,
			},
		},
		security.NewTokenManager(
			"abcdefghijklmnopqrstuvwxyz-1234567890",
			time.Hour,
		),
	)

	result, err := service.Login(
		context.Background(),
		"admin@test.local",
		"Password123!",
	)
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	if result.AccessToken == "" || result.User.Role != models.RoleAdmin {
		t.Fatalf("hasil login tidak sesuai: %+v", result)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	passwordHash, _ := security.HashPassword("Password123!")

	service := NewAuthService(
		fakeUserFinder{
			user: &models.User{
				ID:           1,
				Email:        "admin@test.local",
				PasswordHash: passwordHash,
				Role:         models.RoleAdmin,
				IsActive:     true,
			},
		},
		security.NewTokenManager(
			"abcdefghijklmnopqrstuvwxyz-1234567890",
			time.Hour,
		),
	)

	_, err := service.Login(
		context.Background(),
		"admin@test.local",
		"PasswordSalah",
	)
	if err != ErrInvalidCredentials {
		t.Fatalf("Login() error = %v", err)
	}
}

func TestLoginRejectsMissingUser(t *testing.T) {
	service := NewAuthService(
		fakeUserFinder{err: repositories.ErrNotFound},
		security.NewTokenManager(
			"abcdefghijklmnopqrstuvwxyz-1234567890",
			time.Hour,
		),
	)

	_, err := service.Login(
		context.Background(),
		"missing@test.local",
		"Password123!",
	)
	if err != ErrInvalidCredentials {
		t.Fatalf("Login() error = %v", err)
	}
}
