package handlers

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"jedi-reimbursement-system/backend/internal/dto"
	"jedi-reimbursement-system/backend/internal/middleware"
	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/services"
	"jedi-reimbursement-system/backend/internal/storage"
)

type ClaimReviewHandler struct {
	service  *services.ClaimReviewService
	receipts *storage.ReceiptStorage
}

func NewClaimReviewHandler(
	service *services.ClaimReviewService,
	receipts *storage.ReceiptStorage,
) *ClaimReviewHandler {
	return &ClaimReviewHandler{
		service:  service,
		receipts: receipts,
	}
}

func (h *ClaimReviewHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.List(
		r.Context(),
		r.URL.Query().Get("type"),
		r.URL.Query().Get("status"),
		positiveQueryInt(r, "page", 1),
		positiveQueryInt(r, "limit", 10),
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	items := make(
		[]dto.ClaimReviewResponse,
		0,
		len(result.Items),
	)
	for index := range result.Items {
		items = append(
			items,
			claimReviewResponse(&result.Items[index]),
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Daftar pemeriksaan berhasil dimuat",
			Data: dto.ClaimReviewListResponse{
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

func (h *ClaimReviewHandler) Detail(
	w http.ResponseWriter,
	r *http.Request,
) {
	claimID, err := pathInt64(r, "id")
	if err != nil {
		writeInvalidClaimID(w)
		return
	}

	item, err := h.service.Detail(
		r.Context(),
		r.PathValue("type"),
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
			Message: "Detail pemeriksaan berhasil dimuat",
			Data:    claimReviewResponse(item),
		},
	)
}

func (h *ClaimReviewHandler) Review(
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
		writeInvalidClaimID(w)
		return
	}

	var request dto.ReviewClaimRequest
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

	item, err := h.service.Review(
		r.Context(),
		services.ReviewClaimInput{
			ClaimType:  r.PathValue("type"),
			ClaimID:    claimID,
			Status:     request.Status,
			Note:       request.Note,
			ReviewerID: authUser.ID,
		},
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
			Message: "Status klaim berhasil diperbarui",
			Data:    claimReviewResponse(item),
		},
	)
}

func (h *ClaimReviewHandler) History(
	w http.ResponseWriter,
	r *http.Request,
) {
	claimID, err := pathInt64(r, "id")
	if err != nil {
		writeInvalidClaimID(w)
		return
	}

	history, err := h.service.History(
		r.Context(),
		r.PathValue("type"),
		claimID,
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	items := make(
		[]dto.ClaimHistoryResponse,
		0,
		len(history),
	)
	for index := range history {
		items = append(
			items,
			dto.ClaimHistoryResponse{
				ID:             history[index].ID,
				PreviousStatus: history[index].PreviousStatus,
				NewStatus:      history[index].NewStatus,
				Note:           history[index].Note,
				UpdatedBy:      history[index].UpdatedBy,
				UpdatedByName:  history[index].UpdatedByName,
				CreatedAt: history[index].CreatedAt.Format(
					"2006-01-02T15:04:05Z07:00",
				),
			},
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Riwayat status berhasil dimuat",
			Data: dto.ClaimHistoryListResponse{
				Items: items,
			},
		},
	)
}

func (h *ClaimReviewHandler) Receipt(
	w http.ResponseWriter,
	r *http.Request,
) {
	claimID, err := pathInt64(r, "id")
	if err != nil {
		writeInvalidClaimID(w)
		return
	}

	item, err := h.service.Detail(
		r.Context(),
		r.PathValue("type"),
		claimID,
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	if item.ClaimType != models.ClaimTypeParking ||
		strings.TrimSpace(item.ReceiptPath) == "" {
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Bukti pembayaran tidak tersedia",
			},
		)
		return
	}

	file, err := h.receipts.Open(item.ReceiptPath)
	if err != nil {
		writeReceiptStorageError(w, err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Bukti pembayaran tidak dapat dibaca",
			},
		)
		return
	}

	contentType := item.ReceiptMIMEType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	fileName := item.ReceiptOriginalName
	if fileName == "" {
		fileName = fmt.Sprintf("bukti-parkir-%d", item.ClaimID)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType(
			"inline",
			map[string]string{"filename": fileName},
		),
	)
	w.Header().Set(
		"Content-Length",
		strconv.FormatInt(stat.Size(), 10),
	)
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file)
}

func (h *ClaimReviewHandler) writeServiceError(
	w http.ResponseWriter,
	err error,
) {
	var validationError *services.ClaimReviewValidationError

	switch {
	case errors.As(err, &validationError):
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Data pemeriksaan tidak valid",
				Errors:  validationError.Fields,
			},
		)
	case errors.Is(err, services.ErrClaimFinalized):
		responses.WriteJSON(
			w,
			http.StatusConflict,
			responses.APIResponse{
				Success: false,
				Message: "Klaim ini sudah diputuskan",
			},
		)
	case errors.Is(err, repositories.ErrNotFound):
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Klaim tidak ditemukan",
			},
		)
	case errors.Is(err, services.ErrInvalidClaimReview):
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "Permintaan pemeriksaan tidak valid",
			},
		)
	default:
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Pemeriksaan klaim tidak dapat diproses",
			},
		)
	}
}

func claimReviewResponse(
	item *models.ClaimReview,
) dto.ClaimReviewResponse {
	return dto.ClaimReviewResponse{
		ClaimType:           item.ClaimType,
		ClaimID:             item.ClaimID,
		EmployeeID:          item.EmployeeID,
		EmployeeNumber:      item.EmployeeNumber,
		EmployeeName:        item.EmployeeName,
		ClaimDate:           item.ClaimDate.Format("2006-01-02"),
		ClaimEndDate:        item.ClaimEndDate.Format("2006-01-02"),
		Title:               item.Title,
		Description:         item.Description,
		Amount:              item.Amount,
		StartTime:           item.StartTime,
		EndTime:             item.EndTime,
		DurationHours:       item.DurationHours,
		Status:              item.Status,
		AdminNote:           item.AdminNote,
		ReviewerName:        item.ReviewerName,
		ReviewedAt:          dto.FormatOptionalTime(item.ReviewedAt),
		CreatedAt:           item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:           item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ReceiptAvailable:    item.ReceiptPath != "",
		ReceiptOriginalName: item.ReceiptOriginalName,
		ReceiptMIMEType:     item.ReceiptMIMEType,
		ReceiptSize:         item.ReceiptSize,
	}
}

func writeInvalidClaimID(w http.ResponseWriter) {
	responses.WriteJSON(
		w,
		http.StatusBadRequest,
		responses.APIResponse{
			Success: false,
			Message: "ID klaim tidak valid",
		},
	)
}

func writeReceiptStorageError(
	w http.ResponseWriter,
	err error,
) {
	if errors.Is(err, storage.ErrReceiptUnavailable) {
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Bukti pembayaran tidak tersedia",
			},
		)
		return
	}

	responses.WriteJSON(
		w,
		http.StatusInternalServerError,
		responses.APIResponse{
			Success: false,
			Message: "Bukti pembayaran tidak dapat diproses",
		},
	)
}
