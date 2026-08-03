package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type Migration struct {
	Name  string
	Query string
}

type Migrator struct {
	db *sql.DB
}

func NewMigrator(db *sql.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Up(ctx context.Context) ([]string, error) {
	if err := m.ensureMigrationTable(ctx); err != nil {
		return nil, err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return nil, err
	}

	applied := make([]string, 0)

	for _, migration := range migrations {
		alreadyApplied, err := m.isApplied(ctx, migration.Name)
		if err != nil {
			return applied, err
		}

		if alreadyApplied {
			continue
		}

		// MySQL menjalankan banyak perintah DDL dengan implicit commit.
		// Karena setiap migration menggunakan CREATE TABLE IF NOT EXISTS,
		// eksekusi dapat diulang dengan aman apabila pencatatan gagal.
		if _, err := m.db.ExecContext(ctx, migration.Query); err != nil {
			return applied, fmt.Errorf(
				"menjalankan migration %s: %w",
				migration.Name,
				err,
			)
		}

		if _, err := m.db.ExecContext(
			ctx,
			`INSERT INTO schema_migrations (name, applied_at)
             VALUES (?, ?)`,
			migration.Name,
			time.Now(),
		); err != nil {
			return applied, fmt.Errorf(
				"mencatat migration %s: %w",
				migration.Name,
				err,
			)
		}

		applied = append(applied, migration.Name)
	}

	return applied, nil
}

func (m *Migrator) ensureMigrationTable(ctx context.Context) error {
	const query = `
        CREATE TABLE IF NOT EXISTS schema_migrations (
            id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
            name VARCHAR(255) NOT NULL UNIQUE,
            applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `

	if _, err := m.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("membuat tabel schema_migrations: %w", err)
	}

	return nil
}

func (m *Migrator) isApplied(
	ctx context.Context,
	name string,
) (bool, error) {
	var count int

	err := m.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM schema_migrations WHERE name = ?`,
		name,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf(
			"memeriksa migration %s: %w",
			name,
			err,
		)
	}

	return count > 0, nil
}

func loadMigrations() ([]Migration, error) {
	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return nil, fmt.Errorf("membaca folder migration: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	migrations := make([]Migration, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		query, err := migrationFiles.ReadFile(
			"migrations/" + entry.Name(),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"membaca migration %s: %w",
				entry.Name(),
				err,
			)
		}

		migrations = append(migrations, Migration{
			Name:  entry.Name(),
			Query: string(query),
		})
	}

	return migrations, nil
}
