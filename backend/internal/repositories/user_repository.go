package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"jedi-reimbursement-system/backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *models.User,
) error {
	const query = `
		INSERT INTO users (
			employee_number,
			full_name,
			email,
			password_hash,
			position,
			division,
			role,
			is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		strings.TrimSpace(user.EmployeeNumber),
		strings.TrimSpace(user.FullName),
		strings.ToLower(strings.TrimSpace(user.Email)),
		user.PasswordHash,
		nullableString(user.Position),
		nullableString(user.Division),
		user.Role,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("membuat pengguna: %w", err)
	}

	user.ID, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("membaca ID pengguna: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.User, error) {
	const query = `
		SELECT
			id,
			employee_number,
			full_name,
			email,
			password_hash,
			COALESCE(position, ''),
			COALESCE(division, ''),
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	return scanUser(r.db.QueryRowContext(ctx, query, id))
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	const query = `
		SELECT
			id,
			employee_number,
			full_name,
			email,
			password_hash,
			COALESCE(position, ''),
			COALESCE(division, ''),
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	return scanUser(
		r.db.QueryRowContext(
			ctx,
			query,
			strings.ToLower(strings.TrimSpace(email)),
		),
	)
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	if err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM users`,
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("menghitung pengguna: %w", err)
	}

	return count, nil
}

func scanUser(row *sql.Row) (*models.User, error) {
	user := &models.User{}

	err := row.Scan(
		&user.ID,
		&user.EmployeeNumber,
		&user.FullName,
		&user.Email,
		&user.PasswordHash,
		&user.Position,
		&user.Division,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("membaca pengguna: %w", err)
	}

	return user, nil
}

func nullableString(value string) any {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return value
}
