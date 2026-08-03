package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

const (
	ParkingClaimMaxAmount     = 200_000.0
	ParkingClaimMaxMonthsBack = 3
)

var ErrInvalidParkingClaim = errors.New("data klaim parkir tidak valid")

type ParkingClaimRepository interface {
	Create(
		ctx context.Context,
		claim *models.ParkingClaim,
	) error
	FindByIDAndEmployeeID(
		ctx context.Context,
		id int64,
		employeeID int64,
	) (*models.ParkingClaim, error)
	ListByEmployeeID(
		ctx context.Context,
		employeeID int64,
		limit int,
		offset int,
	) ([]models.ParkingClaim, error)
	CountByEmployeeID(ctx context.Context, employeeID int64) (int64, error)
}

type ParkingClaimValidationError struct {
	Fields map[string]string
}

func (e *ParkingClaimValidationError) Error() string {
	return ErrInvalidParkingClaim.Error()
}

type CreateParkingClaimInput struct {
	EmployeeID          int64
	ParkingStartDate    string
	ParkingEndDate      string
	ParkingLocation     string
	Amount              float64
	Description         string
	ReceiptPath         string
	ReceiptOriginalName string
	ReceiptMIMEType     string
	ReceiptSize         int64
}

type ParkingClaimPage struct {
	Items      []models.ParkingClaim
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type ParkingClaimService struct {
	repository ParkingClaimRepository
	now        func() time.Time
}

func NewParkingClaimService(repository ParkingClaimRepository) *ParkingClaimService {
	return &ParkingClaimService{repository: repository, now: time.Now}
}

func (s *ParkingClaimService) Create(
	ctx context.Context,
	input CreateParkingClaimInput,
) (*models.ParkingClaim, error) {
	startDate, endDate, validationErrors := s.validateCreate(input)
	if len(validationErrors) > 0 {
		return nil, &ParkingClaimValidationError{Fields: validationErrors}
	}

	claim := &models.ParkingClaim{
		EmployeeID:          input.EmployeeID,
		ParkingStartDate:    startDate,
		ParkingEndDate:      endDate,
		ParkingLocation:     strings.TrimSpace(input.ParkingLocation),
		Amount:              input.Amount,
		Description:         strings.TrimSpace(input.Description),
		ReceiptPath:         input.ReceiptPath,
		ReceiptOriginalName: input.ReceiptOriginalName,
		ReceiptMIMEType:     input.ReceiptMIMEType,
		ReceiptSize:         input.ReceiptSize,
		Status:              models.ClaimStatusPending,
	}

	if err := s.repository.Create(ctx, claim); err != nil {
		return nil, fmt.Errorf("menyimpan klaim parkir: %w", err)
	}

	created, err := s.repository.FindByIDAndEmployeeID(
		ctx,
		claim.ID,
		input.EmployeeID,
	)
	if err != nil {
		return nil, fmt.Errorf("memuat klaim parkir: %w", err)
	}
	return created, nil
}

func (s *ParkingClaimService) List(
	ctx context.Context,
	employeeID int64,
	page int,
	limit int,
) (ParkingClaimPage, error) {
	if employeeID < 1 {
		return ParkingClaimPage{}, ErrInvalidParkingClaim
	}
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}
	offset := (page - 1) * limit

	items, err := s.repository.ListByEmployeeID(ctx, employeeID, limit, offset)
	if err != nil {
		return ParkingClaimPage{}, err
	}
	total, err := s.repository.CountByEmployeeID(ctx, employeeID)
	if err != nil {
		return ParkingClaimPage{}, err
	}
	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}
	return ParkingClaimPage{
		Items: items, Page: page, Limit: limit, Total: total, TotalPages: totalPages,
	}, nil
}

func (s *ParkingClaimService) Detail(
	ctx context.Context,
	employeeID int64,
	claimID int64,
) (*models.ParkingClaim, error) {
	if employeeID < 1 || claimID < 1 {
		return nil, ErrInvalidParkingClaim
	}
	return s.repository.FindByIDAndEmployeeID(ctx, claimID, employeeID)
}

func (s *ParkingClaimService) validateCreate(
	input CreateParkingClaimInput,
) (time.Time, time.Time, map[string]string) {
	fields := make(map[string]string)
	location := s.now().Location()

	if input.EmployeeID < 1 {
		fields["employee_id"] = "Pengguna tidak valid"
	}

	startDate, startErr := time.ParseInLocation(
		"2006-01-02", strings.TrimSpace(input.ParkingStartDate), location,
	)
	endDate, endErr := time.ParseInLocation(
		"2006-01-02", strings.TrimSpace(input.ParkingEndDate), location,
	)
	if startErr != nil {
		fields["parking_start_date"] = "Tanggal mulai parkir tidak valid"
	}
	if endErr != nil {
		fields["parking_end_date"] = "Tanggal selesai parkir tidak valid"
	}

	if startErr == nil && endErr == nil {
		today := dateOnly(s.now())
		earliest := firstDayOfMonth(today).AddDate(0, -ParkingClaimMaxMonthsBack, 0)

		if startDate.After(endDate) {
			fields["parking_end_date"] = "Tanggal selesai tidak boleh sebelum tanggal mulai"
		}
		if startDate.Year() != endDate.Year() || startDate.Month() != endDate.Month() {
			fields["parking_end_date"] = "Satu pengajuan hanya boleh mencakup satu bulan kalender"
		}
		if startDate.Before(earliest) {
			fields["parking_start_date"] = "Klaim hanya dapat diajukan untuk bulan berjalan sampai 3 bulan sebelumnya"
		}
		if startDate.After(today) {
			fields["parking_start_date"] = "Tanggal mulai tidak boleh di masa depan"
		}
		if endDate.After(today) {
			fields["parking_end_date"] = "Tanggal selesai tidak boleh di masa depan"
		}
	}

	parkingLocation := strings.TrimSpace(input.ParkingLocation)
	if len(parkingLocation) < 3 {
		fields["parking_location"] = "Lokasi parkir minimal 3 karakter"
	} else if len(parkingLocation) > 200 {
		fields["parking_location"] = "Lokasi parkir maksimal 200 karakter"
	}

	if input.Amount <= 0 {
		fields["amount"] = "Nominal harus lebih besar dari nol"
	} else if input.Amount > ParkingClaimMaxAmount {
		fields["amount"] = "Nominal maksimal setiap pengajuan klaim parkir adalah Rp200.000"
	}

	description := strings.TrimSpace(input.Description)
	if len(description) > 1000 {
		fields["description"] = "Deskripsi maksimal 1000 karakter"
	}

	if strings.TrimSpace(input.ReceiptPath) == "" ||
		strings.TrimSpace(input.ReceiptOriginalName) == "" ||
		strings.TrimSpace(input.ReceiptMIMEType) == "" ||
		input.ReceiptSize < 1 {
		fields["receipt"] = "Bukti pembayaran wajib diunggah"
	}
	return startDate, endDate, fields
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, value.Location())
}

func firstDayOfMonth(value time.Time) time.Time {
	year, month, _ := value.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, value.Location())
}
