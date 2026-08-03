package services

import (
	"context"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

type fakeParkingClaimRepository struct {
	claim models.ParkingClaim
}

func (f *fakeParkingClaimRepository) Create(
	_ context.Context,
	claim *models.ParkingClaim,
) error {
	claim.ID = 10
	claim.CreatedAt = time.Date(2026, 7, 18, 10, 0, 0, 0, time.UTC)
	claim.UpdatedAt = claim.CreatedAt
	f.claim = *claim
	return nil
}

func (f *fakeParkingClaimRepository) FindByIDAndEmployeeID(
	_ context.Context,
	_ int64,
	_ int64,
) (*models.ParkingClaim, error) {
	claim := f.claim
	return &claim, nil
}

func (f *fakeParkingClaimRepository) ListByEmployeeID(
	_ context.Context,
	_ int64,
	_ int,
	_ int,
) ([]models.ParkingClaim, error) {
	return []models.ParkingClaim{f.claim}, nil
}

func (f *fakeParkingClaimRepository) CountByEmployeeID(
	_ context.Context,
	_ int64,
) (int64, error) {
	return 1, nil
}

func validParkingInput() CreateParkingClaimInput {
	return CreateParkingClaimInput{
		EmployeeID:          1,
		ParkingStartDate:    "2026-07-01",
		ParkingEndDate:      "2026-07-18",
		ParkingLocation:     "Gedung Kantor",
		Amount:              20000,
		Description:         "Rekap parkir bulan Juli",
		ReceiptPath:         "parking-receipts/test.pdf",
		ReceiptOriginalName: "bukti.pdf",
		ReceiptMIMEType:     "application/pdf",
		ReceiptSize:         100,
	}
}

func newParkingService(repository ParkingClaimRepository) *ParkingClaimService {
	service := NewParkingClaimService(repository)
	service.now = func() time.Time {
		return time.Date(2026, 7, 30, 12, 0, 0, 0, time.UTC)
	}
	return service
}

func TestCreateParkingClaimWithDateRange(t *testing.T) {
	repository := &fakeParkingClaimRepository{}
	claim, err := newParkingService(repository).Create(
		context.Background(), validParkingInput(),
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if claim.ID != 10 || claim.Status != models.ClaimStatusPending {
		t.Fatalf("claim = %#v", claim)
	}
	if claim.ParkingStartDate.Format("2006-01-02") != "2026-07-01" ||
		claim.ParkingEndDate.Format("2006-01-02") != "2026-07-18" {
		t.Fatal("rentang tanggal tidak tersimpan")
	}
}

func TestCreateRejectsCrossMonthRange(t *testing.T) {
	input := validParkingInput()
	input.ParkingStartDate = "2026-06-30"
	input.ParkingEndDate = "2026-07-01"
	_, err := newParkingService(&fakeParkingClaimRepository{}).Create(
		context.Background(), input,
	)
	validationError, ok := err.(*ParkingClaimValidationError)
	if !ok || validationError.Fields["parking_end_date"] == "" {
		t.Fatalf("error = %T %v", err, err)
	}
}

func TestCreateRejectsClaimOlderThanThreeMonths(t *testing.T) {
	input := validParkingInput()
	input.ParkingStartDate = "2026-03-31"
	input.ParkingEndDate = "2026-03-31"
	_, err := newParkingService(&fakeParkingClaimRepository{}).Create(
		context.Background(), input,
	)
	validationError, ok := err.(*ParkingClaimValidationError)
	if !ok || validationError.Fields["parking_start_date"] == "" {
		t.Fatalf("error = %T %v", err, err)
	}
}

func TestCreateAcceptsFirstDayThreeMonthsAgo(t *testing.T) {
	input := validParkingInput()
	input.ParkingStartDate = "2026-04-01"
	input.ParkingEndDate = "2026-04-30"
	input.Amount = 200000

	_, err := newParkingService(&fakeParkingClaimRepository{}).Create(
		context.Background(), input,
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestCreateAcceptsMaximumAmountPerSubmission(t *testing.T) {
	input := validParkingInput()
	input.Amount = 200000

	_, err := newParkingService(&fakeParkingClaimRepository{}).Create(
		context.Background(), input,
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestCreateRejectsAmountAboveMaximumPerSubmission(t *testing.T) {
	input := validParkingInput()
	input.Amount = 200001
	_, err := newParkingService(&fakeParkingClaimRepository{}).Create(
		context.Background(), input,
	)
	validationError, ok := err.(*ParkingClaimValidationError)
	if !ok || validationError.Fields["amount"] == "" {
		t.Fatalf("error = %T %v", err, err)
	}
}
