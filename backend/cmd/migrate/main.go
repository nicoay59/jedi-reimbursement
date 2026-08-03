package main

import (
	"context"
	"log"
	"time"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/database"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("gagal memuat konfigurasi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg)
	if err != nil {
		log.Fatalf("gagal menghubungkan MySQL: %v", err)
	}
	defer db.Close()

	migrator := database.NewMigrator(db)

	applied, err := migrator.Up(ctx)
	if err != nil {
		log.Fatalf("migration gagal: %v", err)
	}

	if len(applied) == 0 {
		log.Println("tidak ada migration baru")
		return
	}

	for _, name := range applied {
		log.Printf("migration berhasil: %s", name)
	}
}
