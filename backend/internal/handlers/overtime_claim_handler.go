package handlers

import (
	"errors"
	"net/http"

	"jedi-reimbursement-system/backend/internal/dto"
	"jedi-reimbursement-system/backend/internal/middleware"
	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/services"
)

type OvertimeClaimHandler struct {
	service *services.OvertimeClaimService
}

func NewOvertimeClaimHandler(
	service *services.OvertimeClaimService,
) *OvertimeClaimHandler {
	return &OvertimeClaimHandler{service: service}
}

func (h *OvertimeClaimHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	authUser, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	var request dto.CreateOvertimeClaimRequest
	if err := decodeJSON(w, r, &request); err != nil {
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)
		return
	}

	claim, err := h.service.Create(
		r.Context(),
		services.CreateOvertimeClaimInput{
			EmployeeID:      authUser.ID,
			OvertimeDate:    request.OvertimeDate,
			StartTime:       request.StartTime,
			EndTime:         request.EndTime,
			WorkDescription: request.WorkDescription,
		},
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	responses.WriteJSON(
		w,
		http.StatusCreated,
		responses.APIResponse{
			Success: true,
			Message: "Klaim lembur berhasil diajukan",
			Data:    overtimeClaimResponse(claim),
		},
	)
}

func (h *OvertimeClaimHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	authUser, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	page := positiveQueryInt(r, "page", 1)
	limit := positiveQueryInt(r, "limit", 10)

	result, err := h.service.List(
		r.Context(),
		authUser.ID,
		page,
		limit,
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	items := make(
		[]dto.OvertimeClaimResponse,
		0,
		len(result.Items),
	)
	for index := range result.Items {
		items = append(
			items,
			overtimeClaimResponse(&result.Items[index]),
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Daftar klaim lembur berhasil dimuat",
			Data: dto.OvertimeClaimListResponse{
				Items: items,
				Pagination: dto.PaginationResponse{
					Page:       result.Page,
					Limit:      result.Limit,
					Total:      result.Total,
					TotalPages: result.TotalPages,
				},
			},
		},
	)
}

func (h *OvertimeClaimHandler) Detail(
	w http.ResponseWriter,
	r *http.Request,
) {
	authUser, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	claimID, err := pathInt64(r, "id")
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "ID klaim lembur tidak valid",
			},
		)
		return
	}

	claim, err := h.service.Detail(
		r.Context(),
		authUser.ID,
		claimID,
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Detail klaim lembur berhasil dimuat",
			Data:    overtimeClaimResponse(claim),
		},
	)
}

func (h *OvertimeClaimHandler) writeServiceError(
	w http.ResponseWriter,
	err error,
) {
	var validationError *services.OvertimeClaimValidationError

	switch {
	case errors.As(err, &validationError):
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Data klaim lembur tidak valid",
				Errors:  validationError.Fields,
			},
		)
	case errors.Is(err, repositories.ErrNotFound):
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Klaim lembur tidak ditemukan",
			},
		)
	case errors.Is(err, services.ErrInvalidOvertimeClaim):
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "Permintaan klaim lembur tidak valid",
			},
		)
	default:
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Klaim lembur tidak dapat diproses",
			},
		)
	}
}

func overtimeClaimResponse(
	claim *models.OvertimeClaim,
) dto.OvertimeClaimResponse {
	return dto.NewOvertimeClaimResponse(
		claim.ID,
		claim.OvertimeDate,
		claim.StartTime,
		claim.EndTime,
		claim.DurationHours,
		claim.WorkDescription,
		claim.Status,
		claim.AdminNote,
		claim.CreatedAt,
		claim.UpdatedAt,
	)
}
