package services

import (
	"context"
	"math"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

type fakeOvertimeClaimRepository struct {
	claim models.OvertimeClaim
}

func (f *fakeOvertimeClaimRepository) Create(
	_ context.Context,
	claim *models.OvertimeClaim,
) error {
	claim.ID = 20
	claim.CreatedAt = time.Date(
		2026, 7, 18, 20, 0, 0, 0, time.UTC,
	)
	claim.UpdatedAt = claim.CreatedAt
	f.claim = *claim
	return nil
}

func (f *fakeOvertimeClaimRepository) FindByIDAndEmployeeID(
	_ context.Context,
	_ int64,
	_ int64,
) (*models.OvertimeClaim, error) {
	claim := f.claim
	return &claim, nil
}

func (f *fakeOvertimeClaimRepository) ListByEmployeeID(
	_ context.Context,
	_ int64,
	_ int,
	_ int,
) ([]models.OvertimeClaim, error) {
	return []models.OvertimeClaim{f.claim}, nil
}

func (f *fakeOvertimeClaimRepository) CountByEmployeeID(
	_ context.Context,
	_ int64,
) (int64, error) {
	return 1, nil
}

func TestCreateOvertimeClaimSameDay(t *testing.T) {
	repository := &fakeOvertimeClaimRepository{}
	service := NewOvertimeClaimService(repository)
	service.now = func() time.Time {
		return time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	}

	claim, err := service.Create(
		context.Background(),
		CreateOvertimeClaimInput{
			EmployeeID:      1,
			OvertimeDate:    "2026-07-18",
			StartTime:       "18:00",
			EndTime:         "21:30",
			WorkDescription: "Menyelesaikan laporan bulanan",
		},
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if math.Abs(claim.DurationHours-3.5) > 0.001 {
		t.Fatalf("DurationHours = %v", claim.DurationHours)
	}

	if claim.Status != models.ClaimStatusPending {
		t.Fatalf("Status = %q", claim.Status)
	}
}

func TestCreateOvertimeClaimAcrossMidnight(t *testing.T) {
	repository := &fakeOvertimeClaimRepository{}
	service := NewOvertimeClaimService(repository)
	service.now = func() time.Time {
		return time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	}

	claim, err := service.Create(
		context.Background(),
		CreateOvertimeClaimInput{
			EmployeeID:      1,
			OvertimeDate:    "2026-07-18",
			StartTime:       "22:00",
			EndTime:         "01:30",
			WorkDescription: "Pemeliharaan sistem produksi",
		},
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if math.Abs(claim.DurationHours-3.5) > 0.001 {
		t.Fatalf("DurationHours = %v", claim.DurationHours)
	}
}

func TestCreateRejectsExcessiveDuration(t *testing.T) {
	service := NewOvertimeClaimService(
		&fakeOvertimeClaimRepository{},
	)
	service.now = func() time.Time {
		return time.Date(2026, 7, 18, 12, 0, 0, 0, time.UTC)
	}

	_, err := service.Create(
		context.Background(),
		CreateOvertimeClaimInput{
			EmployeeID:      1,
			OvertimeDate:    "2026-07-18",
			StartTime:       "08:00",
			EndTime:         "01:00",
			WorkDescription: "Pemeliharaan sistem produksi",
		},
	)

	validationError, ok := err.(*OvertimeClaimValidationError)
	if !ok {
		t.Fatalf("error = %T %v", err, err)
	}

	if validationError.Fields["end_time"] == "" {
		t.Fatal("end_time seharusnya memiliki error")
	}
}
