# Backend Final 1.2.0

Backend Golang menyediakan autentikasi, klaim parkir, klaim lembur,
persetujuan, dashboard, laporan, dan ekspor CSV.

## Mode Lokal

```cmd
copy .env.example .env
go mod tidy
go test ./...
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/api
```

## Command

```text
cmd/api       REST API
cmd/migrate   Database migration
cmd/seed      Akun awal
cmd/wait-db   Menunggu database untuk Docker
```

## Docker

Docker image backend menjalankan:

```text
wait-db → migrate → seed → api
```

Upload bukti disimpan pada `/app/storage/uploads`.
