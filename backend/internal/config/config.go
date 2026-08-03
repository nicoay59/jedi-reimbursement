package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName     string
	AppEnv      string
	AppHost     string
	AppPort     string
	FrontendURL string

	DBHost            string
	DBPort            string
	DBName            string
	DBUser            string
	DBPassword        string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration

	JWTSecret    string
	JWTExpiresIn time.Duration

	UploadDir              string
	ParkingReceiptMaxBytes int64

	AdminEmployeeNumber string
	AdminFullName       string
	AdminEmail          string
	AdminPassword       string
	AdminPosition       string
	AdminDivision       string

	EmployeeEmployeeNumber string
	EmployeeFullName       string
	EmployeeEmail          string
	EmployeePassword       string
	EmployeePosition       string
	EmployeeDivision       string
}

func Load(envPath string) (Config, error) {
	if err := loadEnvFile(envPath); err != nil &&
		!errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	maxOpenConns, err := getEnvInt("DB_MAX_OPEN_CONNS", 20)
	if err != nil {
		return Config{}, err
	}

	maxIdleConns, err := getEnvInt("DB_MAX_IDLE_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	maxLifetimeMinutes, err := getEnvInt(
		"DB_CONN_MAX_LIFETIME_MINUTES",
		30,
	)
	if err != nil {
		return Config{}, err
	}

	jwtExpiresInMinutes, err := getEnvInt(
		"JWT_EXPIRES_IN_MINUTES",
		480,
	)
	if err != nil {
		return Config{}, err
	}

	receiptMaxMB, err := getEnvInt("PARKING_RECEIPT_MAX_MB", 5)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppName:     getEnv("APP_NAME", "Jedi Reimbursement API"),
		AppEnv:      getEnv("APP_ENV", "development"),
		AppHost:     getEnv("APP_HOST", "0.0.0.0"),
		AppPort:     getEnv("APP_PORT", "8080"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBName:            getEnv("DB_NAME", "jedi_reimbursement"),
		DBUser:            getEnv("DB_USER", "jedi_user"),
		DBPassword:        getEnv("DB_PASSWORD", "jedi_password"),
		DBMaxOpenConns:    maxOpenConns,
		DBMaxIdleConns:    maxIdleConns,
		DBConnMaxLifetime: time.Duration(maxLifetimeMinutes) * time.Minute,

		JWTSecret: getEnv(
			"JWT_SECRET",
			"change-this-jwt-secret-before-production-2026",
		),
		JWTExpiresIn: time.Duration(jwtExpiresInMinutes) * time.Minute,

		UploadDir:              getEnv("UPLOAD_DIR", "storage/uploads"),
		ParkingReceiptMaxBytes: int64(receiptMaxMB) * 1024 * 1024,

		AdminEmployeeNumber: getEnv("ADMIN_EMPLOYEE_NUMBER", "ADM001"),
		AdminFullName:       getEnv("ADMIN_FULL_NAME", "Administrator"),
		AdminEmail:          getEnv("ADMIN_EMAIL", "admin@jedi.local"),
		AdminPassword:       getEnv("ADMIN_PASSWORD", "Admin123!"),
		AdminPosition:       getEnv("ADMIN_POSITION", "System Administrator"),
		AdminDivision:       getEnv("ADMIN_DIVISION", "Information Technology"),

		EmployeeEmployeeNumber: getEnv("EMPLOYEE_EMPLOYEE_NUMBER", "EMP001"),
		EmployeeFullName:       getEnv("EMPLOYEE_FULL_NAME", "Karyawan Contoh"),
		EmployeeEmail:          getEnv("EMPLOYEE_EMAIL", "employee@jedi.local"),
		EmployeePassword:       getEnv("EMPLOYEE_PASSWORD", "Employee123!"),
		EmployeePosition:       getEnv("EMPLOYEE_POSITION", "Staff"),
		EmployeeDivision:       getEnv("EMPLOYEE_DIVISION", "Operations"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Address() string {
	return c.AppHost + ":" + c.AppPort
}

func (c Config) Validate() error {
	required := map[string]string{
		"APP_NAME":       c.AppName,
		"APP_HOST":       c.AppHost,
		"APP_PORT":       c.AppPort,
		"FRONTEND_URL":   c.FrontendURL,
		"DB_HOST":        c.DBHost,
		"DB_PORT":        c.DBPort,
		"DB_NAME":        c.DBName,
		"DB_USER":        c.DBUser,
		"JWT_SECRET":     c.JWTSecret,
		"UPLOAD_DIR":     c.UploadDir,
		"ADMIN_EMAIL":    c.AdminEmail,
		"EMPLOYEE_EMAIL": c.EmployeeEmail,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s tidak boleh kosong", key)
		}
	}

	if c.DBMaxOpenConns < 1 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS minimal 1")
	}

	if c.DBMaxIdleConns < 0 || c.DBMaxIdleConns > c.DBMaxOpenConns {
		return fmt.Errorf("nilai DB_MAX_IDLE_CONNS tidak valid")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET minimal 32 karakter")
	}

	if c.JWTExpiresIn < time.Minute {
		return fmt.Errorf("JWT_EXPIRES_IN_MINUTES minimal 1")
	}

	if c.ParkingReceiptMaxBytes < 1024 {
		return fmt.Errorf("PARKING_RECEIPT_MAX_MB minimal 1")
	}

	if len(c.AdminPassword) < 8 || len(c.EmployeePassword) < 8 {
		return fmt.Errorf("password seed minimal 8 karakter")
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) (int, error) {
	raw := getEnv(key, strconv.Itoa(fallback))
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s harus berupa angka: %w", key, err)
	}
	return value, nil
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return fmt.Errorf("gagal mengatur %s: %w", key, err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("membaca file environment: %w", err)
	}

	return nil
}
