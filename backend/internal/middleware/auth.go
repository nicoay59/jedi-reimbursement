package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/security"
)

type authUserKey struct{}

type AuthUser struct {
	ID       int64
	Role     string
	Email    string
	FullName string
}

func Auth(
	tokens *security.TokenManager,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			writeUnauthorized(w, "Token autentikasi diperlukan")
			return
		}

		claims, err := tokens.Parse(token)
		if errors.Is(err, security.ErrExpiredToken) {
			writeUnauthorized(w, "Sesi login telah berakhir")
			return
		}
		if err != nil {
			writeUnauthorized(w, "Token autentikasi tidak valid")
			return
		}

		user := AuthUser{
			ID:       claims.Subject,
			Role:     claims.Role,
			Email:    claims.Email,
			FullName: claims.FullName,
		}

		ctx := context.WithValue(r.Context(), authUserKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthUserFromContext(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(authUserKey{}).(AuthUser)
	return user, ok
}

func bearerToken(value string) string {
	scheme, token, found := strings.Cut(strings.TrimSpace(value), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return ""
	}

	return strings.TrimSpace(token)
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	responses.WriteJSON(
		w,
		http.StatusUnauthorized,
		responses.APIResponse{
			Success: false,
			Message: message,
		},
	)
}
