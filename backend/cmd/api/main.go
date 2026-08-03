package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/database"
	"jedi-reimbursement-system/backend/internal/routes"
	"jedi-reimbursement-system/backend/internal/storage"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("gagal memuat konfigurasi: %v", err)
	}

	startupContext, cancelStartup := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancelStartup()

	db, err := database.Open(startupContext, cfg)
	if err != nil {
		log.Fatalf("gagal menghubungkan MySQL: %v", err)
	}
	defer db.Close()

	receiptStorage, err := storage.NewReceiptStorage(
		cfg.UploadDir,
		cfg.ParkingReceiptMaxBytes,
	)
	if err != nil {
		log.Fatalf("gagal menyiapkan penyimpanan bukti: %v", err)
	}

	server := &http.Server{
		Addr: cfg.Address(),
		Handler: routes.New(
			cfg,
			db,
			receiptStorage,
		),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf(
			"%s berjalan di http://localhost:%s",
			cfg.AppName,
			cfg.AppPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("server berhenti karena kesalahan: %v", err)
	case signalValue := <-signals:
		log.Printf("menerima sinyal %s", signalValue)
	}

	shutdownContext, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownContext); err != nil {
		log.Printf("server gagal dihentikan secara normal: %v", err)
	}
}
