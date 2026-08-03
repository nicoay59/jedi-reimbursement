package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

const claimReviewSelect = `
	SELECT
		'PARKING' AS claim_type,
		parking.id AS claim_id,
		parking.employee_id,
		employee.employee_number,
		employee.full_name AS employee_name,
		parking.parking_date AS claim_date,
		COALESCE(parking.parking_end_date, parking.parking_date) AS claim_end_date,
		parking.parking_location AS title,
		COALESCE(parking.description, '') AS description,
		parking.amount,
		'' AS start_time,
		'' AS end_time,
		0 AS duration_hours,
		parking.status,
		COALESCE(parking.admin_note, '') AS admin_note,
		parking.reviewed_by,
		COALESCE(reviewer.full_name, '') AS reviewer_name,
		parking.reviewed_at,
		parking.created_at,
		parking.updated_at,
		COALESCE(parking.receipt_path, '') AS receipt_path,
		COALESCE(parking.receipt_original_name, '') AS receipt_original_name,
		COALESCE(parking.receipt_mime_type, '') AS receipt_mime_type,
		COALESCE(parking.receipt_size, 0) AS receipt_size
	FROM parking_claims AS parking
	INNER JOIN users AS employee
		ON employee.id = parking.employee_id
	LEFT JOIN users AS reviewer
		ON reviewer.id = parking.reviewed_by

	UNION ALL

	SELECT
		'OVERTIME' AS claim_type,
		overtime.id AS claim_id,
		overtime.employee_id,
		employee.employee_number,
		employee.full_name AS employee_name,
		overtime.overtime_date AS claim_date,
		overtime.overtime_date AS claim_end_date,
		'Klaim Lembur' AS title,
		overtime.work_description AS description,
		0 AS amount,
		CAST(overtime.start_time AS CHAR) AS start_time,
		CAST(overtime.end_time AS CHAR) AS end_time,
		overtime.duration_hours,
		overtime.status,
		COALESCE(overtime.admin_note, '') AS admin_note,
		overtime.reviewed_by,
		COALESCE(reviewer.full_name, '') AS reviewer_name,
		overtime.reviewed_at,
		overtime.created_at,
		overtime.updated_at,
		'' AS receipt_path,
		'' AS receipt_original_name,
		'' AS receipt_mime_type,
		0 AS receipt_size
	FROM overtime_claims AS overtime
	INNER JOIN users AS employee
		ON employee.id = overtime.employee_id
	LEFT JOIN users AS reviewer
		ON reviewer.id = overtime.reviewed_by
`

type ClaimReviewRepository struct {
	db *sql.DB
}

func NewClaimReviewRepository(
	db *sql.DB,
) *ClaimReviewRepository {
	return &ClaimReviewRepository{db: db}
}

func (r *ClaimReviewRepository) List(
	ctx context.Context,
	claimType string,
	status string,
	limit int,
	offset int,
) ([]models.ClaimReview, error) {
	query := `
		SELECT *
		FROM (` + claimReviewSelect + `) AS claims
		WHERE (? = 'ALL' OR claim_type = ?)
		  AND (? = 'ALL' OR status = ?)
		ORDER BY created_at DESC, claim_id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		claimType,
		claimType,
		status,
		status,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar pemeriksaan: %w", err)
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
		return nil, fmt.Errorf("membaca daftar pemeriksaan: %w", err)
	}

	return items, nil
}

func (r *ClaimReviewRepository) Count(
	ctx context.Context,
	claimType string,
	status string,
) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM (` + claimReviewSelect + `) AS claims
		WHERE (? = 'ALL' OR claim_type = ?)
		  AND (? = 'ALL' OR status = ?)
	`

	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		query,
		claimType,
		claimType,
		status,
		status,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("menghitung daftar pemeriksaan: %w", err)
	}

	return total, nil
}

func (r *ClaimReviewRepository) FindByTypeAndID(
	ctx context.Context,
	claimType string,
	claimID int64,
) (*models.ClaimReview, error) {
	query := `
		SELECT *
		FROM (` + claimReviewSelect + `) AS claims
		WHERE claim_type = ? AND claim_id = ?
		LIMIT 1
	`

	return scanClaimReview(
		r.db.QueryRowContext(ctx, query, claimType, claimID),
	)
}

func (r *ClaimReviewRepository) UpdateStatus(
	ctx context.Context,
	claimType string,
	claimID int64,
	newStatus string,
	note string,
	reviewerID int64,
) error {
	tableName, err := reviewTableName(claimType)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("memulai transaksi pemeriksaan: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var previousStatus string
	queryStatus := fmt.Sprintf(
		"SELECT status FROM %s WHERE id = ? FOR UPDATE",
		tableName,
	)

	err = tx.QueryRowContext(
		ctx,
		queryStatus,
		claimID,
	).Scan(&previousStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("membaca status klaim: %w", err)
	}

	if previousStatus != models.ClaimStatusPending {
		return ErrConflict
	}

	reviewedAt := time.Now()
	updateQuery := fmt.Sprintf(
		`UPDATE %s
		 SET status = ?,
		     admin_note = ?,
		     reviewed_by = ?,
		     reviewed_at = ?,
		     updated_at = ?
		 WHERE id = ?`,
		tableName,
	)

	if _, err := tx.ExecContext(
		ctx,
		updateQuery,
		newStatus,
		nullableString(note),
		reviewerID,
		reviewedAt,
		reviewedAt,
		claimID,
	); err != nil {
		return fmt.Errorf("memperbarui status klaim: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO claim_status_histories (
			claim_type,
			claim_id,
			previous_status,
			new_status,
			note,
			updated_by,
			created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		claimType,
		claimID,
		previousStatus,
		newStatus,
		nullableString(note),
		reviewerID,
		reviewedAt,
	); err != nil {
		return fmt.Errorf("mencatat riwayat status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("menyelesaikan transaksi pemeriksaan: %w", err)
	}

	committed = true
	return nil
}

func (r *ClaimReviewRepository) History(
	ctx context.Context,
	claimType string,
	claimID int64,
) ([]models.ClaimReviewHistory, error) {
	const query = `
		SELECT
			history.id,
			history.claim_type,
			history.claim_id,
			COALESCE(history.previous_status, ''),
			history.new_status,
			COALESCE(history.note, ''),
			history.updated_by,
			user.full_name,
			history.created_at
		FROM claim_status_histories AS history
		INNER JOIN users AS user
			ON user.id = history.updated_by
		WHERE history.claim_type = ?
		  AND history.claim_id = ?
		ORDER BY history.created_at ASC, history.id ASC
	`

	rows, err := r.db.QueryContext(ctx, query, claimType, claimID)
	if err != nil {
		return nil, fmt.Errorf("mengambil riwayat status: %w", err)
	}
	defer rows.Close()

	items := make([]models.ClaimReviewHistory, 0)

	for rows.Next() {
		item := models.ClaimReviewHistory{}

		if err := rows.Scan(
			&item.ID,
			&item.ClaimType,
			&item.ClaimID,
			&item.PreviousStatus,
			&item.NewStatus,
			&item.Note,
			&item.UpdatedBy,
			&item.UpdatedByName,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("membaca riwayat status: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca riwayat status: %w", err)
	}

	return items, nil
}

func scanClaimReview(
	row rowScanner,
) (*models.ClaimReview, error) {
	item := &models.ClaimReview{}
	var reviewedBy sql.NullInt64
	var reviewedAt sql.NullTime

	err := row.Scan(
		&item.ClaimType,
		&item.ClaimID,
		&item.EmployeeID,
		&item.EmployeeNumber,
		&item.EmployeeName,
		&item.ClaimDate,
		&item.ClaimEndDate,
		&item.Title,
		&item.Description,
		&item.Amount,
		&item.StartTime,
		&item.EndTime,
		&item.DurationHours,
		&item.Status,
		&item.AdminNote,
		&reviewedBy,
		&item.ReviewerName,
		&reviewedAt,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.ReceiptPath,
		&item.ReceiptOriginalName,
		&item.ReceiptMIMEType,
		&item.ReceiptSize,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("membaca data pemeriksaan: %w", err)
	}

	if reviewedBy.Valid {
		value := reviewedBy.Int64
		item.ReviewedBy = &value
	}

	if reviewedAt.Valid {
		value := reviewedAt.Time
		item.ReviewedAt = &value
	}

	item.StartTime = normalizeDatabaseClock(item.StartTime)
	item.EndTime = normalizeDatabaseClock(item.EndTime)

	return item, nil
}

func reviewTableName(claimType string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(claimType)) {
	case models.ClaimTypeParking:
		return "parking_claims", nil
	case models.ClaimTypeOvertime:
		return "overtime_claims", nil
	default:
		return "", ErrNotFound
	}
}

func normalizeDatabaseClock(value string) string {
	if len(value) >= 5 {
		return value[:5]
	}
	return value
}

func insertInitialHistory(
	ctx context.Context,
	tx *sql.Tx,
	claimType string,
	claimID int64,
	employeeID int64,
	status string,
) error {
	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO claim_status_histories (
			claim_type,
			claim_id,
			previous_status,
			new_status,
			note,
			updated_by
		) VALUES (?, ?, NULL, ?, ?, ?)`,
		claimType,
		claimID,
		status,
		"Pengajuan dibuat",
		employeeID,
	); err != nil {
		return fmt.Errorf("mencatat status awal: %w", err)
	}

	return nil
}
