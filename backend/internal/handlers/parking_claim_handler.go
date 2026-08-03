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

type ParkingClaimHandler struct {
	service      *services.ParkingClaimService
	receipts     *storage.ReceiptStorage
	maxBodyBytes int64
}

func NewParkingClaimHandler(
	service *services.ParkingClaimService,
	receipts *storage.ReceiptStorage,
	maxReceiptBytes int64,
) *ParkingClaimHandler {
	return &ParkingClaimHandler{
		service:      service,
		receipts:     receipts,
		maxBodyBytes: maxReceiptBytes + 1024*1024,
	}
}

func (h *ParkingClaimHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	authUser, ok := middleware.AuthUserFromContext(r.Context())
	if !ok {
		writeUnauthorized(w)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.maxBodyBytes)

	if err := r.ParseMultipartForm(h.maxBodyBytes); err != nil {
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "Form tidak valid atau ukuran request terlalu besar",
			},
		)
		return
	}

	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	amount, err := strconv.ParseFloat(
		strings.TrimSpace(r.FormValue("amount")),
		64,
	)
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Data klaim parkir tidak valid",
				Errors: map[string]string{
					"amount": "Nominal harus berupa angka",
				},
			},
		)
		return
	}

	file, fileHeader, err := r.FormFile("receipt")
	if err != nil {
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Data klaim parkir tidak valid",
				Errors: map[string]string{
					"receipt": "Bukti pembayaran wajib diunggah",
				},
			},
		)
		return
	}
	defer file.Close()

	savedReceipt, err := h.receipts.Save(file, fileHeader)
	if err != nil {
		h.writeStorageError(w, err)
		return
	}

	claim, err := h.service.Create(
		r.Context(),
		services.CreateParkingClaimInput{
			EmployeeID:          authUser.ID,
			ParkingStartDate:    r.FormValue("parking_start_date"),
			ParkingEndDate:      r.FormValue("parking_end_date"),
			ParkingLocation:     r.FormValue("parking_location"),
			Amount:              amount,
			Description:         r.FormValue("description"),
			ReceiptPath:         savedReceipt.RelativePath,
			ReceiptOriginalName: savedReceipt.OriginalName,
			ReceiptMIMEType:     savedReceipt.MIMEType,
			ReceiptSize:         savedReceipt.Size,
		},
	)
	if err != nil {
		_ = h.receipts.Delete(savedReceipt.RelativePath)
		h.writeServiceError(w, err)
		return
	}

	responses.WriteJSON(
		w,
		http.StatusCreated,
		responses.APIResponse{
			Success: true,
			Message: "Klaim parkir berhasil diajukan",
			Data:    parkingClaimResponse(claim),
		},
	)
}

func (h *ParkingClaimHandler) List(
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
		[]dto.ParkingClaimResponse,
		0,
		len(result.Items),
	)
	for index := range result.Items {
		items = append(
			items,
			parkingClaimResponse(&result.Items[index]),
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Daftar klaim parkir berhasil dimuat",
			Data: dto.ParkingClaimListResponse{
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

func (h *ParkingClaimHandler) Detail(
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
				Message: "ID klaim parkir tidak valid",
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
			Message: "Detail klaim parkir berhasil dimuat",
			Data:    parkingClaimResponse(claim),
		},
	)
}

func (h *ParkingClaimHandler) Receipt(
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
				Message: "ID klaim parkir tidak valid",
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

	file, err := h.receipts.Open(claim.ReceiptPath)
	if err != nil {
		h.writeStorageError(w, err)
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

	contentType := claim.ReceiptMIMEType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	fileName := claim.ReceiptOriginalName
	if fileName == "" {
		fileName = fmt.Sprintf("bukti-parkir-%d", claim.ID)
	}

	disposition := mime.FormatMediaType(
		"inline",
		map[string]string{"filename": fileName},
	)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", disposition)
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, file); err != nil {
		return
	}
}

func (h *ParkingClaimHandler) writeServiceError(
	w http.ResponseWriter,
	err error,
) {
	var validationError *services.ParkingClaimValidationError

	switch {
	case errors.As(err, &validationError):
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Data klaim parkir tidak valid",
				Errors:  validationError.Fields,
			},
		)
	case errors.Is(err, repositories.ErrNotFound):
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Klaim parkir tidak ditemukan",
			},
		)
	case errors.Is(err, services.ErrInvalidParkingClaim):
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "Permintaan klaim parkir tidak valid",
			},
		)
	default:
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Klaim parkir tidak dapat diproses",
			},
		)
	}
}

func (h *ParkingClaimHandler) writeStorageError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(err, storage.ErrReceiptTooLarge):
		responses.WriteJSON(
			w,
			http.StatusRequestEntityTooLarge,
			responses.APIResponse{
				Success: false,
				Message: "Ukuran bukti pembayaran terlalu besar",
				Errors: map[string]string{
					"receipt": "Ukuran maksimal bukti adalah 5 MB",
				},
			},
		)
	case errors.Is(err, storage.ErrReceiptType):
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Jenis bukti pembayaran tidak didukung",
				Errors: map[string]string{
					"receipt": "Gunakan file JPG, PNG, atau PDF",
				},
			},
		)
	case errors.Is(err, storage.ErrReceiptUnavailable):
		responses.WriteJSON(
			w,
			http.StatusNotFound,
			responses.APIResponse{
				Success: false,
				Message: "Bukti pembayaran tidak tersedia",
			},
		)
	default:
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Bukti pembayaran tidak dapat diproses",
			},
		)
	}
}

func parkingClaimResponse(
	claim *models.ParkingClaim,
) dto.ParkingClaimResponse {
	return dto.NewParkingClaimResponse(
		claim.ID,
		claim.ParkingStartDate,
		claim.ParkingEndDate,
		claim.ParkingLocation,
		claim.Amount,
		claim.Description,
		claim.ReceiptOriginalName,
		claim.ReceiptMIMEType,
		claim.ReceiptSize,
		claim.Status,
		claim.AdminNote,
		claim.CreatedAt,
		claim.UpdatedAt,
	)
}

func positiveQueryInt(
	r *http.Request,
	key string,
	fallback int,
) int {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil || number < 1 {
		return fallback
	}

	return number
}

func pathInt64(
	r *http.Request,
	key string,
) (int64, error) {
	value := strings.TrimSpace(r.PathValue(key))
	number, err := strconv.ParseInt(value, 10, 64)
	if err != nil || number < 1 {
		return 0, errors.New("ID tidak valid")
	}

	return number, nil
}

func writeUnauthorized(w http.ResponseWriter) {
	responses.WriteJSON(
		w,
		http.StatusUnauthorized,
		responses.APIResponse{
			Success: false,
			Message: "Token autentikasi diperlukan",
		},
	)
}
