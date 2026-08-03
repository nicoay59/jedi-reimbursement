package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"jedi-reimbursement-system/backend/internal/models"
)

type ParkingClaimRepository struct {
	db *sql.DB
}

func NewParkingClaimRepository(
	db *sql.DB,
) *ParkingClaimRepository {
	return &ParkingClaimRepository{db: db}
}

func (r *ParkingClaimRepository) Create(
	ctx context.Context,
	claim *models.ParkingClaim,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("memulai transaksi klaim parkir: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	const query = `
		INSERT INTO parking_claims (
			employee_id,
			parking_date,
			parking_end_date,
			parking_location,
			amount,
			description,
			receipt_path,
			receipt_original_name,
			receipt_mime_type,
			receipt_size,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.ExecContext(
		ctx,
		query,
		claim.EmployeeID,
		claim.ParkingStartDate,
		claim.ParkingEndDate,
		claim.ParkingLocation,
		claim.Amount,
		nullableString(claim.Description),
		claim.ReceiptPath,
		claim.ReceiptOriginalName,
		claim.ReceiptMIMEType,
		claim.ReceiptSize,
		claim.Status,
	)
	if err != nil {
		return fmt.Errorf("membuat klaim parkir: %w", err)
	}

	claim.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("membaca ID klaim parkir: %w", err)
	}

	if err := insertInitialHistory(
		ctx,
		tx,
		models.ClaimTypeParking,
		claim.ID,
		claim.EmployeeID,
		claim.Status,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("menyelesaikan transaksi klaim parkir: %w", err)
	}

	committed = true
	return nil
}

func (r *ParkingClaimRepository) FindByIDAndEmployeeID(
	ctx context.Context,
	id int64,
	employeeID int64,
) (*models.ParkingClaim, error) {
	const query = `
		SELECT
			id,
			employee_id,
			parking_date,
			COALESCE(parking_end_date, parking_date),
			parking_location,
			amount,
			COALESCE(description, ''),
			COALESCE(receipt_path, ''),
			COALESCE(receipt_original_name, ''),
			COALESCE(receipt_mime_type, ''),
			COALESCE(receipt_size, 0),
			status,
			COALESCE(admin_note, ''),
			reviewed_by,
			reviewed_at,
			created_at,
			updated_at
		FROM parking_claims
		WHERE id = ? AND employee_id = ?
		LIMIT 1
	`

	return scanParkingClaim(
		r.db.QueryRowContext(ctx, query, id, employeeID),
	)
}

func (r *ParkingClaimRepository) ListByEmployeeID(
	ctx context.Context,
	employeeID int64,
	limit int,
	offset int,
) ([]models.ParkingClaim, error) {
	const query = `
		SELECT
			id,
			employee_id,
			parking_date,
			COALESCE(parking_end_date, parking_date),
			parking_location,
			amount,
			COALESCE(description, ''),
			COALESCE(receipt_path, ''),
			COALESCE(receipt_original_name, ''),
			COALESCE(receipt_mime_type, ''),
			COALESCE(receipt_size, 0),
			status,
			COALESCE(admin_note, ''),
			reviewed_by,
			reviewed_at,
			created_at,
			updated_at
		FROM parking_claims
		WHERE employee_id = ?
		ORDER BY parking_date DESC, created_at DESC, id DESC
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
		return nil, fmt.Errorf("mengambil daftar klaim parkir: %w", err)
	}
	defer rows.Close()

	claims := make([]models.ParkingClaim, 0)
	for rows.Next() {
		claim, err := scanParkingClaimRows(rows)
		if err != nil {
			return nil, err
		}
		claims = append(claims, *claim)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca daftar klaim parkir: %w", err)
	}
	return claims, nil
}

func (r *ParkingClaimRepository) CountByEmployeeID(
	ctx context.Context,
	employeeID int64,
) (int64, error) {
	var total int64
	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM parking_claims WHERE employee_id = ?`,
		employeeID,
	).Scan(&total); err != nil {
		return 0, fmt.Errorf("menghitung klaim parkir: %w", err)
	}
	return total, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanParkingClaim(row rowScanner) (*models.ParkingClaim, error) {
	claim := &models.ParkingClaim{}
	var reviewedBy sql.NullInt64
	var reviewedAt sql.NullTime

	err := row.Scan(
		&claim.ID,
		&claim.EmployeeID,
		&claim.ParkingStartDate,
		&claim.ParkingEndDate,
		&claim.ParkingLocation,
		&claim.Amount,
		&claim.Description,
		&claim.ReceiptPath,
		&claim.ReceiptOriginalName,
		&claim.ReceiptMIMEType,
		&claim.ReceiptSize,
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
		return nil, fmt.Errorf("membaca klaim parkir: %w", err)
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

func scanParkingClaimRows(rows *sql.Rows) (*models.ParkingClaim, error) {
	return scanParkingClaim(rows)
}
