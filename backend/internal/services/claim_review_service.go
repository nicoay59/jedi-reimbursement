package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
)

var (
	ErrInvalidClaimReview = errors.New("permintaan pemeriksaan tidak valid")
	ErrClaimFinalized     = errors.New("klaim sudah diputuskan")
)

type ClaimReviewRepository interface {
	List(
		ctx context.Context,
		claimType string,
		status string,
		limit int,
		offset int,
	) ([]models.ClaimReview, error)
	Count(
		ctx context.Context,
		claimType string,
		status string,
	) (int64, error)
	FindByTypeAndID(
		ctx context.Context,
		claimType string,
		claimID int64,
	) (*models.ClaimReview, error)
	UpdateStatus(
		ctx context.Context,
		claimType string,
		claimID int64,
		newStatus string,
		note string,
		reviewerID int64,
	) error
	History(
		ctx context.Context,
		claimType string,
		claimID int64,
	) ([]models.ClaimReviewHistory, error)
}

type ClaimReviewValidationError struct {
	Fields map[string]string
}

func (e *ClaimReviewValidationError) Error() string {
	return ErrInvalidClaimReview.Error()
}

type ClaimReviewPage struct {
	Items      []models.ClaimReview
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type ReviewClaimInput struct {
	ClaimType  string
	ClaimID    int64
	Status     string
	Note       string
	ReviewerID int64
}

type ClaimReviewService struct {
	repository ClaimReviewRepository
}

func NewClaimReviewService(
	repository ClaimReviewRepository,
) *ClaimReviewService {
	return &ClaimReviewService{repository: repository}
}

func (s *ClaimReviewService) List(
	ctx context.Context,
	claimType string,
	status string,
	page int,
	limit int,
) (ClaimReviewPage, error) {
	normalizedType, err := normalizeClaimType(claimType, true)
	if err != nil {
		return ClaimReviewPage{}, err
	}

	normalizedStatus, err := normalizeClaimStatus(status, true)
	if err != nil {
		return ClaimReviewPage{}, err
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

	items, err := s.repository.List(
		ctx,
		normalizedType,
		normalizedStatus,
		limit,
		offset,
	)
	if err != nil {
		return ClaimReviewPage{}, err
	}

	total, err := s.repository.Count(
		ctx,
		normalizedType,
		normalizedStatus,
	)
	if err != nil {
		return ClaimReviewPage{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return ClaimReviewPage{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *ClaimReviewService) Detail(
	ctx context.Context,
	claimType string,
	claimID int64,
) (*models.ClaimReview, error) {
	normalizedType, err := normalizeClaimType(claimType, false)
	if err != nil || claimID < 1 {
		return nil, ErrInvalidClaimReview
	}

	return s.repository.FindByTypeAndID(
		ctx,
		normalizedType,
		claimID,
	)
}

func (s *ClaimReviewService) Review(
	ctx context.Context,
	input ReviewClaimInput,
) (*models.ClaimReview, error) {
	normalizedType, err := normalizeClaimType(
		input.ClaimType,
		false,
	)
	if err != nil {
		return nil, ErrInvalidClaimReview
	}

	if input.ClaimID < 1 || input.ReviewerID < 1 {
		return nil, ErrInvalidClaimReview
	}

	status, err := normalizeReviewDecision(input.Status)
	if err != nil {
		return nil, &ClaimReviewValidationError{
			Fields: map[string]string{
				"status": "Status harus APPROVED atau REJECTED",
			},
		}
	}

	note := strings.TrimSpace(input.Note)
	fields := make(map[string]string)

	if status == models.ClaimStatusRejected && len(note) < 5 {
		fields["note"] =
			"Catatan penolakan minimal 5 karakter"
	}

	if len(note) > 1000 {
		fields["note"] = "Catatan maksimal 1000 karakter"
	}

	if len(fields) > 0 {
		return nil, &ClaimReviewValidationError{
			Fields: fields,
		}
	}

	current, err := s.repository.FindByTypeAndID(
		ctx,
		normalizedType,
		input.ClaimID,
	)
	if err != nil {
		return nil, err
	}

	if current.Status != models.ClaimStatusPending {
		return nil, ErrClaimFinalized
	}

	err = s.repository.UpdateStatus(
		ctx,
		normalizedType,
		input.ClaimID,
		status,
		note,
		input.ReviewerID,
	)
	if errors.Is(err, repositories.ErrConflict) {
		return nil, ErrClaimFinalized
	}
	if err != nil {
		return nil, fmt.Errorf("memperbarui pemeriksaan: %w", err)
	}

	return s.repository.FindByTypeAndID(
		ctx,
		normalizedType,
		input.ClaimID,
	)
}

func (s *ClaimReviewService) History(
	ctx context.Context,
	claimType string,
	claimID int64,
) ([]models.ClaimReviewHistory, error) {
	normalizedType, err := normalizeClaimType(claimType, false)
	if err != nil || claimID < 1 {
		return nil, ErrInvalidClaimReview
	}

	if _, err := s.repository.FindByTypeAndID(
		ctx,
		normalizedType,
		claimID,
	); err != nil {
		return nil, err
	}

	return s.repository.History(
		ctx,
		normalizedType,
		claimID,
	)
}

func normalizeClaimType(
	value string,
	allowAll bool,
) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" && allowAll {
		return models.ClaimTypeAll, nil
	}

	switch value {
	case models.ClaimTypeParking, models.ClaimTypeOvertime:
		return value, nil
	case models.ClaimTypeAll:
		if allowAll {
			return value, nil
		}
	}

	return "", ErrInvalidClaimReview
}

func normalizeClaimStatus(
	value string,
	allowAll bool,
) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" && allowAll {
		return models.ClaimTypeAll, nil
	}

	switch value {
	case models.ClaimStatusPending,
		models.ClaimStatusApproved,
		models.ClaimStatusRejected:
		return value, nil
	case models.ClaimTypeAll:
		if allowAll {
			return value, nil
		}
	}

	return "", ErrInvalidClaimReview
}

func normalizeReviewDecision(
	value string,
) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))

	switch value {
	case models.ClaimStatusApproved,
		models.ClaimStatusRejected:
		return value, nil
	default:
		return "", ErrInvalidClaimReview
	}
}
