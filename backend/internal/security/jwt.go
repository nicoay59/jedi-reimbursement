package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token telah kedaluwarsa")
)

type Claims struct {
	Subject   int64  `json:"sub"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
		now:    time.Now,
	}
}

func (m *TokenManager) Generate(user *models.User) (string, error) {
	now := m.now().UTC()

	headerJSON, err := json.Marshal(map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("membuat header token: %w", err)
	}

	claimsJSON, err := json.Marshal(Claims{
		Subject:   user.ID,
		Role:      user.Role,
		Email:     user.Email,
		FullName:  user.FullName,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("membuat payload token: %w", err)
	}

	headerPart := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsPart := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsignedToken := headerPart + "." + claimsPart
	signature := base64.RawURLEncoding.EncodeToString(m.sign(unsignedToken))

	return unsignedToken + "." + signature, nil
}

func (m *TokenManager) Parse(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	providedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	expectedSignature := m.sign(unsignedToken)
	if len(providedSignature) != len(expectedSignature) ||
		subtle.ConstantTimeCompare(providedSignature, expectedSignature) != 1 {
		return Claims{}, ErrInvalidToken
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil ||
		header.Algorithm != "HS256" ||
		header.Type != "JWT" {
		return Claims{}, ErrInvalidToken
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}

	if claims.Subject < 1 ||
		strings.TrimSpace(claims.Role) == "" ||
		claims.ExpiresAt < 1 {
		return Claims{}, ErrInvalidToken
	}

	if m.now().UTC().Unix() >= claims.ExpiresAt {
		return Claims{}, ErrExpiredToken
	}

	return claims, nil
}

func (m *TokenManager) TTLSeconds() int64 {
	return int64(m.ttl.Seconds())
}

func (m *TokenManager) sign(value string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	_, _ = mac.Write([]byte(value))
	return mac.Sum(nil)
}
