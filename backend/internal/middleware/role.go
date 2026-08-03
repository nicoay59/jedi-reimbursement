package middleware

import (
	"net/http"

	"jedi-reimbursement-system/backend/internal/responses"
)

func RequireRoles(
	allowedRoles ...string,
) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := AuthUserFromContext(r.Context())
			if !ok {
				writeUnauthorized(w, "Token autentikasi diperlukan")
				return
			}

			if _, exists := allowed[user.Role]; !exists {
				responses.WriteJSON(
					w,
					http.StatusForbidden,
					responses.APIResponse{
						Success: false,
						Message: "Anda tidak memiliki hak akses",
					},
				)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
