package services

import (
	"context"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
)

type fakeClaimReviewRepository struct {
	item      models.ClaimReview
	history   []models.ClaimReviewHistory
	updateErr error
}

func (f *fakeClaimReviewRepository) List(
	context.Context,
	string,
	string,
	int,
	int,
) ([]models.ClaimReview, error) {
	return []models.ClaimReview{f.item}, nil
}

func (f *fakeClaimReviewRepository) Count(
	context.Context,
	string,
	string,
) (int64, error) {
	return 1, nil
}

func (f *fakeClaimReviewRepository) FindByTypeAndID(
	context.Context,
	string,
	int64,
) (*models.ClaimReview, error) {
	item := f.item
	return &item, nil
}

func (f *fakeClaimReviewRepository) UpdateStatus(
	_ context.Context,
	_ string,
	_ int64,
	status string,
	note string,
	reviewerID int64,
) error {
	if f.updateErr != nil {
		return f.updateErr
	}

	f.item.Status = status
	f.item.AdminNote = note
	f.item.ReviewedBy = &reviewerID
	now := time.Date(2026, 7, 18, 14, 0, 0, 0, time.UTC)
	f.item.ReviewedAt = &now
	return nil
}

func (f *fakeClaimReviewRepository) History(
	context.Context,
	string,
	int64,
) ([]models.ClaimReviewHistory, error) {
	return f.history, nil
}

func TestReviewApprove(t *testing.T) {
	repository := &fakeClaimReviewRepository{
		item: models.ClaimReview{
			ClaimType: models.ClaimTypeParking,
			ClaimID:   1,
			Status:    models.ClaimStatusPending,
		},
	}

	service := NewClaimReviewService(repository)

	result, err := service.Review(
		context.Background(),
		ReviewClaimInput{
			ClaimType:  models.ClaimTypeParking,
			ClaimID:    1,
			Status:     models.ClaimStatusApproved,
			Note:       "Bukti telah sesuai",
			ReviewerID: 2,
		},
	)
	if err != nil {
		t.Fatalf("Review() error = %v", err)
	}

	if result.Status != models.ClaimStatusApproved {
		t.Fatalf("Status = %q", result.Status)
	}
}

func TestReviewRejectRequiresNote(t *testing.T) {
	service := NewClaimReviewService(
		&fakeClaimReviewRepository{
			item: models.ClaimReview{
				ClaimType: models.ClaimTypeOvertime,
				ClaimID:   2,
				Status:    models.ClaimStatusPending,
			},
		},
	)

	_, err := service.Review(
		context.Background(),
		ReviewClaimInput{
			ClaimType:  models.ClaimTypeOvertime,
			ClaimID:    2,
			Status:     models.ClaimStatusRejected,
			Note:       "",
			ReviewerID: 1,
		},
	)

	validationError, ok := err.(*ClaimReviewValidationError)
	if !ok {
		t.Fatalf("error = %T %v", err, err)
	}

	if validationError.Fields["note"] == "" {
		t.Fatal("note seharusnya memiliki error")
	}
}

func TestReviewFinalizedClaim(t *testing.T) {
	service := NewClaimReviewService(
		&fakeClaimReviewRepository{
			item: models.ClaimReview{
				ClaimType: models.ClaimTypeParking,
				ClaimID:   1,
				Status:    models.ClaimStatusApproved,
			},
		},
	)

	_, err := service.Review(
		context.Background(),
		ReviewClaimInput{
			ClaimType:  models.ClaimTypeParking,
			ClaimID:    1,
			Status:     models.ClaimStatusRejected,
			Note:       "Data tidak sesuai",
			ReviewerID: 2,
		},
	)
	if err != ErrClaimFinalized {
		t.Fatalf("Review() error = %v", err)
	}
}

func TestReviewHandlesRepositoryConflict(t *testing.T) {
	service := NewClaimReviewService(
		&fakeClaimReviewRepository{
			item: models.ClaimReview{
				ClaimType: models.ClaimTypeParking,
				ClaimID:   1,
				Status:    models.ClaimStatusPending,
			},
			updateErr: repositories.ErrConflict,
		},
	)

	_, err := service.Review(
		context.Background(),
		ReviewClaimInput{
			ClaimType:  models.ClaimTypeParking,
			ClaimID:    1,
			Status:     models.ClaimStatusApproved,
			ReviewerID: 2,
		},
	)
	if err != ErrClaimFinalized {
		t.Fatalf("Review() error = %v", err)
	}
}
