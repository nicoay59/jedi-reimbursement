package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"

	"jedi-reimbursement-system/backend/internal/config"
)

type Pinger interface {
	PingContext(ctx context.Context) error
}

func Open(ctx context.Context, cfg config.Config) (*sql.DB, error) {
	driverConfig := mysqlDriver.NewConfig()
	driverConfig.User = cfg.DBUser
	driverConfig.Passwd = cfg.DBPassword
	driverConfig.Net = "tcp"
	driverConfig.Addr = net.JoinHostPort(cfg.DBHost, cfg.DBPort)
	driverConfig.DBName = cfg.DBName
	driverConfig.ParseTime = true
	driverConfig.Loc = time.Local
	driverConfig.Params = map[string]string{
		"charset":   "utf8mb4",
		"collation": "utf8mb4_unicode_ci",
	}

	db, err := sql.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("membuka database: %w", err)
	}

	db.SetMaxOpenConns(cfg.DBMaxOpenConns)
	db.SetMaxIdleConns(cfg.DBMaxIdleConns)
	db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("memeriksa koneksi database: %w", err)
	}

	return db, nil
}
