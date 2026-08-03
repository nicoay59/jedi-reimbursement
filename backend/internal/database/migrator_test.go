package database

import "testing"

func TestLoadMigrations(t *testing.T) {
	migrations, err := loadMigrations()
	if err != nil {
		t.Fatalf("loadMigrations() error = %v", err)
	}

	if len(migrations) != 7 {
		t.Fatalf("jumlah migration = %d, ingin 7", len(migrations))
	}

	expectedFirst := "001_create_users.sql"
	if migrations[0].Name != expectedFirst {
		t.Fatalf(
			"migration pertama = %q, ingin %q",
			migrations[0].Name,
			expectedFirst,
		)
	}
}
