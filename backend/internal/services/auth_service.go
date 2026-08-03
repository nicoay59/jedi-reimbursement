package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"jedi-reimbursement-system/backend/internal/dto"
	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/security"
)

var (
	ErrInvalidCredentials = errors.New("email atau password salah")
	ErrInactiveUser       = errors.New("akun pengguna tidak aktif")
)

type UserFinder interface {
	FindByEmail(context.Context, string) (*models.User, error)
	FindByID(context.Context, int64) (*models.User, error)
}

type AuthService struct {
	users  UserFinder
	tokens *security.TokenManager
}

func NewAuthService(
	users UserFinder,
	tokens *security.TokenManager,
) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (dto.LoginResponse, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, repositories.ErrNotFound) {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("mencari pengguna: %w", err)
	}

	if !user.IsActive {
		return dto.LoginResponse{}, ErrInactiveUser
	}

	if !security.VerifyPassword(user.PasswordHash, password) {
		return dto.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("membuat access token: %w", err)
	}

	return dto.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.tokens.TTLSeconds(),
		User:        mapUserResponse(user),
	}, nil
}

func (s *AuthService) Me(
	ctx context.Context,
	userID int64,
) (dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if !user.IsActive {
		return dto.UserResponse{}, ErrInactiveUser
	}

	return mapUserResponse(user), nil
}

func mapUserResponse(user *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:             user.ID,
		EmployeeNumber: user.EmployeeNumber,
		FullName:       user.FullName,
		Email:          user.Email,
		Position:       user.Position,
		Division:       user.Division,
		Role:           user.Role,
	}
}
