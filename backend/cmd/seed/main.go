package main

import (
	"context"
	"errors"
	"log"
	"time"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/database"
	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/security"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("gagal memuat konfigurasi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg)
	if err != nil {
		log.Fatalf("gagal menghubungkan MySQL: %v", err)
	}
	defer db.Close()

	repository := repositories.NewUserRepository(db)

	items := []seedUser{
		{
			label: "administrator",
			user: models.User{
				EmployeeNumber: cfg.AdminEmployeeNumber,
				FullName:       cfg.AdminFullName,
				Email:          cfg.AdminEmail,
				Position:       cfg.AdminPosition,
				Division:       cfg.AdminDivision,
				Role:           models.RoleAdmin,
				IsActive:       true,
			},
			password: cfg.AdminPassword,
		},
		{
			label: "karyawan",
			user: models.User{
				EmployeeNumber: cfg.EmployeeEmployeeNumber,
				FullName:       cfg.EmployeeFullName,
				Email:          cfg.EmployeeEmail,
				Position:       cfg.EmployeePosition,
				Division:       cfg.EmployeeDivision,
				Role:           models.RoleEmployee,
				IsActive:       true,
			},
			password: cfg.EmployeePassword,
		},
	}

	for index := range items {
		if err := createUserIfMissing(
			ctx,
			repository,
			&items[index],
		); err != nil {
			log.Fatalf("seed %s gagal: %v", items[index].label, err)
		}
	}
}

type seedUser struct {
	label    string
	user     models.User
	password string
}

func createUserIfMissing(
	ctx context.Context,
	repository *repositories.UserRepository,
	item *seedUser,
) error {
	existing, err := repository.FindByEmail(ctx, item.user.Email)
	if err == nil && existing != nil {
		log.Printf("%s %s sudah tersedia", item.label, item.user.Email)
		return nil
	}

	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return err
	}

	passwordHash, err := security.HashPassword(item.password)
	if err != nil {
		return err
	}

	item.user.PasswordHash = passwordHash

	if err := repository.Create(ctx, &item.user); err != nil {
		return err
	}

	log.Printf(
		"%s berhasil dibuat: %s (%s)",
		item.label,
		item.user.FullName,
		item.user.Email,
	)

	return nil
}
