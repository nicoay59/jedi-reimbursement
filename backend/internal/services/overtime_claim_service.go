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

var ErrInvalidOvertimeClaim = errors.New("data klaim lembur tidak valid")

type OvertimeClaimRepository interface {
	Create(
		ctx context.Context,
		claim *models.OvertimeClaim,
	) error
	FindByIDAndEmployeeID(
		ctx context.Context,
		id int64,
		employeeID int64,
	) (*models.OvertimeClaim, error)
	ListByEmployeeID(
		ctx context.Context,
		employeeID int64,
		limit int,
		offset int,
	) ([]models.OvertimeClaim, error)
	CountByEmployeeID(
		ctx context.Context,
		employeeID int64,
	) (int64, error)
}

type OvertimeClaimValidationError struct {
	Fields map[string]string
}

func (e *OvertimeClaimValidationError) Error() string {
	return ErrInvalidOvertimeClaim.Error()
}

type CreateOvertimeClaimInput struct {
	EmployeeID      int64
	OvertimeDate    string
	StartTime       string
	EndTime         string
	WorkDescription string
}

type OvertimeClaimPage struct {
	Items      []models.OvertimeClaim
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type OvertimeClaimService struct {
	repository OvertimeClaimRepository
	now        func() time.Time
}

func NewOvertimeClaimService(
	repository OvertimeClaimRepository,
) *OvertimeClaimService {
	return &OvertimeClaimService{
		repository: repository,
		now:        time.Now,
	}
}

func (s *OvertimeClaimService) Create(
	ctx context.Context,
	input CreateOvertimeClaimInput,
) (*models.OvertimeClaim, error) {
	overtimeDate, startTime, endTime, durationHours, validationErrors :=
		s.validateCreate(input)

	if len(validationErrors) > 0 {
		return nil, &OvertimeClaimValidationError{
			Fields: validationErrors,
		}
	}

	claim := &models.OvertimeClaim{
		EmployeeID:      input.EmployeeID,
		OvertimeDate:    overtimeDate,
		StartTime:       startTime,
		EndTime:         endTime,
		DurationHours:   durationHours,
		WorkDescription: strings.TrimSpace(input.WorkDescription),
		Status:          models.ClaimStatusPending,
	}

	if err := s.repository.Create(ctx, claim); err != nil {
		return nil, fmt.Errorf("menyimpan klaim lembur: %w", err)
	}

	created, err := s.repository.FindByIDAndEmployeeID(
		ctx,
		claim.ID,
		input.EmployeeID,
	)
	if err != nil {
		return nil, fmt.Errorf("memuat klaim lembur: %w", err)
	}

	return created, nil
}

func (s *OvertimeClaimService) List(
	ctx context.Context,
	employeeID int64,
	page int,
	limit int,
) (OvertimeClaimPage, error) {
	if employeeID < 1 {
		return OvertimeClaimPage{}, ErrInvalidOvertimeClaim
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

	items, err := s.repository.ListByEmployeeID(
		ctx,
		employeeID,
		limit,
		offset,
	)
	if err != nil {
		return OvertimeClaimPage{}, err
	}

	total, err := s.repository.CountByEmployeeID(ctx, employeeID)
	if err != nil {
		return OvertimeClaimPage{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return OvertimeClaimPage{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *OvertimeClaimService) Detail(
	ctx context.Context,
	employeeID int64,
	claimID int64,
) (*models.OvertimeClaim, error) {
	if employeeID < 1 || claimID < 1 {
		return nil, ErrInvalidOvertimeClaim
	}

	return s.repository.FindByIDAndEmployeeID(
		ctx,
		claimID,
		employeeID,
	)
}

func (s *OvertimeClaimService) validateCreate(
	input CreateOvertimeClaimInput,
) (
	time.Time,
	string,
	string,
	float64,
	map[string]string,
) {
	fields := make(map[string]string)

	if input.EmployeeID < 1 {
		fields["employee_id"] = "Pengguna tidak valid"
	}

	overtimeDate, err := time.ParseInLocation("2006-01-02", input.OvertimeDate, s.now().Location())
	if err != nil {
		fields["overtime_date"] = "Tanggal lembur tidak valid"
	} else {
		today := dateOnly(s.now())
		if overtimeDate.After(today) {
			fields["overtime_date"] =
				"Tanggal lembur tidak boleh di masa depan"
		}
	}

	startMinutes, normalizedStart, err := parseClock(input.StartTime)
	if err != nil {
		fields["start_time"] = "Waktu mulai tidak valid"
	}

	endMinutes, normalizedEnd, err := parseClock(input.EndTime)
	if err != nil {
		fields["end_time"] = "Waktu selesai tidak valid"
	}

	durationMinutes := 0
	if fields["start_time"] == "" && fields["end_time"] == "" {
		durationMinutes = endMinutes - startMinutes

		// Waktu selesai yang sama atau lebih kecil dianggap hari berikutnya.
		if durationMinutes <= 0 {
			durationMinutes += 24 * 60
		}

		if durationMinutes < 30 {
			fields["end_time"] = "Durasi lembur minimal 30 menit"
		}

		if durationMinutes > 16*60 {
			fields["end_time"] = "Durasi lembur maksimal 16 jam"
		}
	}

	description := strings.TrimSpace(input.WorkDescription)
	if len(description) < 10 {
		fields["work_description"] =
			"Deskripsi pekerjaan minimal 10 karakter"
	} else if len(description) > 2000 {
		fields["work_description"] =
			"Deskripsi pekerjaan maksimal 2000 karakter"
	}

	durationHours := math.Round(
		(float64(durationMinutes)/60)*100,
	) / 100

	return overtimeDate,
		normalizedStart,
		normalizedEnd,
		durationHours,
		fields
}

func parseClock(value string) (int, string, error) {
	value = strings.TrimSpace(value)

	parsed, err := time.Parse("15:04", value)
	if err != nil {
		return 0, "", err
	}

	minutes := parsed.Hour()*60 + parsed.Minute()
	return minutes, parsed.Format("15:04:05"), nil
}
