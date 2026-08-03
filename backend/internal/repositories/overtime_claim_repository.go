package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"jedi-reimbursement-system/backend/internal/models"
)

type OvertimeClaimRepository struct {
	db *sql.DB
}

func NewOvertimeClaimRepository(
	db *sql.DB,
) *OvertimeClaimRepository {
	return &OvertimeClaimRepository{db: db}
}

func (r *OvertimeClaimRepository) Create(
	ctx context.Context,
	claim *models.OvertimeClaim,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("memulai transaksi klaim lembur: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	const query = `
		INSERT INTO overtime_claims (
			employee_id,
			overtime_date,
			start_time,
			end_time,
			duration_hours,
			work_description,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		claim.EmployeeID,
		claim.OvertimeDate,
		claim.StartTime,
		claim.EndTime,
		claim.DurationHours,
		claim.WorkDescription,
		claim.Status,
	)
	if err != nil {
		return fmt.Errorf("membuat klaim lembur: %w", err)
	}

	claim.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("membaca ID klaim lembur: %w", err)
	}

	if err := insertInitialHistory(
		ctx,
		tx,
		models.ClaimTypeOvertime,
		claim.ID,
		claim.EmployeeID,
		claim.Status,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("menyelesaikan transaksi klaim lembur: %w", err)
	}

	committed = true
	return nil
}

func (r *OvertimeClaimRepository) FindByIDAndEmployeeID(
	ctx context.Context,
	id int64,
	employeeID int64,
) (*models.OvertimeClaim, error) {
	const query = `
		SELECT
			id,
			employee_id,
			overtime_date,
			CAST(start_time AS CHAR),
			CAST(end_time AS CHAR),
			duration_hours,
			work_description,
			status,
			COALESCE(admin_note, ''),
			reviewed_by,
			reviewed_at,
			created_at,
			updated_at
		FROM overtime_claims
		WHERE id = ? AND employee_id = ?
		LIMIT 1
	`

	return scanOvertimeClaim(
		r.db.QueryRowContext(ctx, query, id, employeeID),
	)
}

func (r *OvertimeClaimRepository) ListByEmployeeID(
	ctx context.Context,
	employeeID int64,
	limit int,
	offset int,
) ([]models.OvertimeClaim, error) {
	const query = `
		SELECT
			id,
			employee_id,
			overtime_date,
			CAST(start_time AS CHAR),
			CAST(end_time AS CHAR),
			duration_hours,
			work_description,
			status,
			COALESCE(admin_note, ''),
			reviewed_by,
			reviewed_at,
			created_at,
			updated_at
		FROM overtime_claims
		WHERE employee_id = ?
		ORDER BY created_at DESC, id DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		employeeID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar klaim lembur: %w", err)
	}
	defer rows.Close()

	claims := make([]models.OvertimeClaim, 0)

	for rows.Next() {
		claim, err := scanOvertimeClaim(rows)
		if err != nil {
			return nil, err
		}

		claims = append(claims, *claim)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca daftar klaim lembur: %w", err)
	}

	return claims, nil
}

func (r *OvertimeClaimRepository) CountByEmployeeID(
	ctx context.Context,
	employeeID int64,
) (int64, error) {
	var total int64

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM overtime_claims WHERE employee_id = ?`,
		employeeID,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("menghitung klaim lembur: %w", err)
	}

	return total, nil
}

func scanOvertimeClaim(
	row rowScanner,
) (*models.OvertimeClaim, error) {
	claim := &models.OvertimeClaim{}
	var reviewedBy sql.NullInt64
	var reviewedAt sql.NullTime

	err := row.Scan(
		&claim.ID,
		&claim.EmployeeID,
		&claim.OvertimeDate,
		&claim.StartTime,
		&claim.EndTime,
		&claim.DurationHours,
		&claim.WorkDescription,
		&claim.Status,
		&claim.AdminNote,
		&reviewedBy,
		&reviewedAt,
		&claim.CreatedAt,
		&claim.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("membaca klaim lembur: %w", err)
	}

	if reviewedBy.Valid {
		value := reviewedBy.Int64
		claim.ReviewedBy = &value
	}

	if reviewedAt.Valid {
		value := reviewedAt.Time
		claim.ReviewedAt = &value
	}

	return claim, nil
}
