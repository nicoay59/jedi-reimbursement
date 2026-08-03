package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"jedi-reimbursement-system/backend/internal/dto"
	"jedi-reimbursement-system/backend/internal/middleware"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/services"
)

const maxJSONBodySize = 1 << 20

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var request dto.LoginRequest
	if err := decodeJSON(w, r, &request); err != nil {
		responses.WriteJSON(w, http.StatusBadRequest, responses.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if strings.TrimSpace(request.Email) == "" ||
		strings.TrimSpace(request.Password) == "" {
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Email dan password wajib diisi",
			},
		)
		return
	}

	result, err := h.service.Login(
		r.Context(),
		request.Email,
		request.Password,
	)
	if errors.Is(err, services.ErrInvalidCredentials) {
		responses.WriteJSON(w, http.StatusUnauthorized, responses.APIResponse{
			Success: false,
			Message: "Email atau password salah",
		})
		return
	}
	if errors.Is(err, services.ErrInactiveUser) {
		responses.WriteJSON(w, http.StatusForbidden, responses.APIResponse{
			Success: false,
			Message: "Akun pengguna tidak aktif",
		})
		return
	}
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Login tidak dapat diproses",
			},
		)
		return
	}

	responses.WriteJSON(w, http.StatusOK, responses.APIResponse{
		Success: true,
		Message: "Login berhasil",
		Data:    result,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, http.StatusUnauthorized, responses.APIResponse{
			Success: false,
			Message: "Token autentikasi diperlukan",
		})
		return
	}

	user, err := h.service.Me(r.Context(), authUser.ID)
	if errors.Is(err, repositories.ErrNotFound) ||
		errors.Is(err, services.ErrInactiveUser) {
		responses.WriteJSON(w, http.StatusUnauthorized, responses.APIResponse{
			Success: false,
			Message: "Sesi pengguna tidak valid",
		})
		return
	}
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Profil pengguna tidak dapat dimuat",
			},
		)
		return
	}

	responses.WriteJSON(w, http.StatusOK, responses.APIResponse{
		Success: true,
		Message: "Profil pengguna berhasil dimuat",
		Data:    user,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, _ *http.Request) {
	responses.WriteJSON(w, http.StatusOK, responses.APIResponse{
		Success: true,
		Message: "Logout berhasil",
	})
}

func decodeJSON(
	w http.ResponseWriter,
	r *http.Request,
	destination any,
) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodySize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		if errors.Is(err, io.EOF) {
			return errors.New("Body request tidak boleh kosong")
		}
		return errors.New("Format JSON tidak valid")
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("Body request hanya boleh berisi satu objek JSON")
	}

	return nil
}
