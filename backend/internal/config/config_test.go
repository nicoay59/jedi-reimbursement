package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesEnvironmentFile(t *testing.T) {
	tempDir := t.TempDir()
	envPath := filepath.Join(tempDir, ".env")

	content := []byte(
		"APP_NAME=Test API\n" +
			"APP_ENV=test\n" +
			"APP_HOST=127.0.0.1\n" +
			"APP_PORT=9090\n" +
			"FRONTEND_URL=http://localhost:4173\n" +
			"DB_HOST=localhost\n" +
			"DB_PORT=3306\n" +
			"DB_NAME=test_database\n" +
			"DB_USER=test_user\n" +
			"JWT_SECRET=abcdefghijklmnopqrstuvwxyz-1234567890\n" +
			"JWT_EXPIRES_IN_MINUTES=60\n" +
			"UPLOAD_DIR=test-uploads\n" +
			"PARKING_RECEIPT_MAX_MB=7\n" +
			"ADMIN_EMAIL=admin@test.local\n" +
			"ADMIN_PASSWORD=Password123!\n" +
			"EMPLOYEE_EMAIL=employee@test.local\n" +
			"EMPLOYEE_PASSWORD=Password123!\n",
	)

	if err := os.WriteFile(envPath, content, 0o600); err != nil {
		t.Fatalf("menulis file env: %v", err)
	}

	keys := []string{
		"APP_NAME", "APP_ENV", "APP_HOST", "APP_PORT", "FRONTEND_URL",
		"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "JWT_SECRET",
		"JWT_EXPIRES_IN_MINUTES", "UPLOAD_DIR", "PARKING_RECEIPT_MAX_MB",
		"ADMIN_EMAIL", "ADMIN_PASSWORD", "EMPLOYEE_EMAIL",
		"EMPLOYEE_PASSWORD",
	}

	previousValues := make(map[string]string, len(keys))
	previousExists := make(map[string]bool, len(keys))

	for _, key := range keys {
		value, exists := os.LookupEnv(key)
		previousValues[key] = value
		previousExists[key] = exists
		_ = os.Unsetenv(key)
	}

	t.Cleanup(func() {
		for _, key := range keys {
			if previousExists[key] {
				_ = os.Setenv(key, previousValues[key])
			} else {
				_ = os.Unsetenv(key)
			}
		}
	})

	cfg, err := Load(envPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.AppName != "Test API" || cfg.DBName != "test_database" {
		t.Fatalf("konfigurasi tidak sesuai: %+v", cfg)
	}

	if cfg.ParkingReceiptMaxBytes != 7*1024*1024 {
		t.Fatalf(
			"ParkingReceiptMaxBytes = %d",
			cfg.ParkingReceiptMaxBytes,
		)
	}
}
