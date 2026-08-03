package handlers

import (
	"net/http"

	"jedi-reimbursement-system/backend/internal/middleware"
	"jedi-reimbursement-system/backend/internal/responses"
)

type AccessHandler struct{}

func NewAccessHandler() *AccessHandler {
	return &AccessHandler{}
}

func (h *AccessHandler) Admin(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.AuthUserFromContext(r.Context())
	responses.WriteJSON(w, http.StatusOK, responses.APIResponse{
		Success: true,
		Message: "Akses administrator berhasil",
		Data: map[string]any{
			"user_id":   user.ID,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}

func (h *AccessHandler) Employee(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.AuthUserFromContext(r.Context())
	responses.WriteJSON(w, http.StatusOK, responses.APIResponse{
		Success: true,
		Message: "Akses karyawan berhasil",
		Data: map[string]any{
			"user_id":   user.ID,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}
