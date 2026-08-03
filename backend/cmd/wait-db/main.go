package main

import (
	"context"
	"log"
	"time"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/database"
)

const (
	maximumWait = 90 * time.Second
	retryDelay  = 3 * time.Second
	pingTimeout = 5 * time.Second
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("gagal memuat konfigurasi: %v", err)
	}

	deadline := time.Now().Add(maximumWait)
	attempt := 0

	for {
		attempt++

		ctx, cancel := context.WithTimeout(
			context.Background(),
			pingTimeout,
		)

		db, err := database.Open(ctx, cfg)
		cancel()

		if err == nil {
			_ = db.Close()
			log.Printf(
				"database siap setelah %d percobaan",
				attempt,
			)
			return
		}

		if time.Now().After(deadline) {
			log.Fatalf(
				"database belum siap setelah %s: %v",
				maximumWait,
				err,
			)
		}

		log.Printf(
			"database belum siap, percobaan %d: %v",
			attempt,
			err,
		)

		time.Sleep(retryDelay)
	}
}
