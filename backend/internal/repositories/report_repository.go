package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"jedi-reimbursement-system/backend/internal/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(
	db *sql.DB,
) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Summary(
	ctx context.Context,
	period models.ReportPeriod,
) (models.DashboardSummary, error) {
	query := `
		SELECT
			COUNT(*) AS total_claims,
			COALESCE(SUM(
				CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END
			), 0) AS pending_claims,
			COALESCE(SUM(
				CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END
			), 0) AS approved_claims,
			COALESCE(SUM(
				CASE WHEN status = 'REJECTED' THEN 1 ELSE 0 END
			), 0) AS rejected_claims,
			COALESCE(SUM(
				CASE WHEN claim_type = 'PARKING' THEN 1 ELSE 0 END
			), 0) AS parking_claims,
			COALESCE(SUM(
				CASE WHEN claim_type = 'OVERTIME' THEN 1 ELSE 0 END
			), 0) AS overtime_claims,
			COALESCE(SUM(
				CASE WHEN claim_type = 'PARKING' THEN amount ELSE 0 END
			), 0) AS total_parking_amount,
			COALESCE(SUM(
				CASE
					WHEN claim_type = 'PARKING'
					 AND status = 'APPROVED'
					THEN amount
					ELSE 0
				END
			), 0) AS approved_parking_amount,
			COALESCE(SUM(
				CASE
					WHEN claim_type = 'OVERTIME'
					THEN duration_hours
					ELSE 0
				END
			), 0) AS total_overtime_hours,
			COALESCE(SUM(
				CASE
					WHEN claim_type = 'OVERTIME'
					 AND status = 'APPROVED'
					THEN duration_hours
					ELSE 0
				END
			), 0) AS approved_overtime_hours
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
	`

	var summary models.DashboardSummary

	if err := r.db.QueryRowContext(
		ctx,
		query,
		period.StartDate,
		period.EndDate,
	).Scan(
		&summary.TotalClaims,
		&summary.PendingClaims,
		&summary.ApprovedClaims,
		&summary.RejectedClaims,
		&summary.ParkingClaims,
		&summary.OvertimeClaims,
		&summary.TotalParkingAmount,
		&summary.ApprovedParkingAmount,
		&summary.TotalOvertimeHours,
		&summary.ApprovedOvertimeHours,
	); err != nil {
		return models.DashboardSummary{},
			fmt.Errorf("mengambil ringkasan dashboard: %w", err)
	}

	return summary, nil
}

func (r *ReportRepository) Trend(
	ctx context.Context,
	period models.ReportPeriod,
) ([]models.DashboardTrend, error) {
	query := `
		SELECT
			claim_date,
			COUNT(*) AS total_claims,
			COALESCE(SUM(
				CASE WHEN claim_type = 'PARKING' THEN 1 ELSE 0 END
			), 0) AS parking_claims,
			COALESCE(SUM(
				CASE WHEN claim_type = 'OVERTIME' THEN 1 ELSE 0 END
			), 0) AS overtime_claims,
			COALESCE(SUM(
				CASE WHEN status = 'APPROVED' THEN 1 ELSE 0 END
			), 0) AS approved_claims
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
		GROUP BY claim_date
		ORDER BY claim_date ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		period.StartDate,
		period.EndDate,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil tren dashboard: %w", err)
	}
	defer rows.Close()

	items := make([]models.DashboardTrend, 0)

	for rows.Next() {
		item := models.DashboardTrend{}

		if err := rows.Scan(
			&item.Date,
			&item.TotalClaims,
			&item.ParkingClaims,
			&item.OvertimeClaims,
			&item.ApprovedClaims,
		); err != nil {
			return nil, fmt.Errorf("membaca tren dashboard: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca tren dashboard: %w", err)
	}

	return items, nil
}

func (r *ReportRepository) Recent(
	ctx context.Context,
	period models.ReportPeriod,
	limit int,
) ([]models.ClaimReview, error) {
	query := `
		SELECT *
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
		ORDER BY created_at DESC, claim_id DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		period.StartDate,
		period.EndDate,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil klaim terbaru: %w", err)
	}
	defer rows.Close()

	items := make([]models.ClaimReview, 0)

	for rows.Next() {
		item, err := scanClaimReview(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca klaim terbaru: %w", err)
	}

	return items, nil
}

func (r *ReportRepository) List(
	ctx context.Context,
	filter models.ReportFilter,
	limit int,
	offset int,
) ([]models.ClaimReview, error) {
	query := `
		SELECT *
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
		  AND (? = 'ALL' OR claim_type = ?)
		  AND (? = 'ALL' OR status = ?)
		ORDER BY claim_date DESC, created_at DESC, claim_id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		filter.Period.StartDate,
		filter.Period.EndDate,
		filter.ClaimType,
		filter.ClaimType,
		filter.Status,
		filter.Status,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil laporan klaim: %w", err)
	}
	defer rows.Close()

	items := make([]models.ClaimReview, 0)

	for rows.Next() {
		item, err := scanClaimReview(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca laporan klaim: %w", err)
	}

	return items, nil
}

func (r *ReportRepository) Count(
	ctx context.Context,
	filter models.ReportFilter,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
		  AND (? = 'ALL' OR claim_type = ?)
		  AND (? = 'ALL' OR status = ?)
	`

	var total int64

	if err := r.db.QueryRowContext(
		ctx,
		query,
		filter.Period.StartDate,
		filter.Period.EndDate,
		filter.ClaimType,
		filter.ClaimType,
		filter.Status,
		filter.Status,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("menghitung laporan klaim: %w", err)
	}

	return total, nil
}

func (r *ReportRepository) Export(
	ctx context.Context,
	filter models.ReportFilter,
) ([]models.ClaimReview, error) {
	query := `
		SELECT *
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_date BETWEEN ? AND ?
		  AND (? = 'ALL' OR claim_type = ?)
		  AND (? = 'ALL' OR status = ?)
		ORDER BY claim_date ASC, created_at ASC, claim_id ASC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		filter.Period.StartDate,
		filter.Period.EndDate,
		filter.ClaimType,
		filter.ClaimType,
		filter.Status,
		filter.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil data ekspor: %w", err)
	}
	defer rows.Close()

	items := make([]models.ClaimReview, 0)

	for rows.Next() {
		item, err := scanClaimReview(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca data ekspor: %w", err)
	}

	return items, nil
}
