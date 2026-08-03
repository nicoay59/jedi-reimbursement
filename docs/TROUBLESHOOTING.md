# Troubleshooting Final

## Full Docker Tidak Dapat Dibuka

Jalankan:

```cmd
docker-logs.bat
```

Periksa status:

```cmd
docker compose --env-file .env.docker ps
```

Service yang benar berstatus:

```text
healthy
```

## Container Backend Restart Berulang

Kemungkinan:

- Database belum sehat.
- Migration gagal.
- Environment tidak valid.
- `JWT_SECRET` kurang dari 32 karakter.
- Password seed kurang dari 8 karakter.

Lihat log:

```cmd
docker logs jedi-reimbursement-backend
```

## Port 8088 Sudah Digunakan

Ubah `.env.docker`:

```text
FRONTEND_PORT=8090
```

Kemudian restart Docker.

## Port 3306 Sudah Digunakan

Ubah:

```text
MYSQL_PORT=3307
```

Catatan: backend Docker tetap menggunakan port internal 3306.

## Mode Pengembangan Tidak Berjalan

Pastikan terminal versi sebelumnya sudah ditutup.

Pemeriksaan manual:

```cmd
netstat -ano | findstr :8080
netstat -ano | findstr :5173
```

## Login Gagal Setelah Mengubah Password Environment

Seed tidak mengganti password pengguna yang sudah ada. Pilihan:

- Gunakan database baru.
- Perbarui password melalui database.
- Jalankan `docker-reset.bat` untuk data demo baru.

`docker-reset.bat` menghapus seluruh data.

## Bukti Parkir Hilang

Pada Full Docker, bukti berada pada volume:

```text
jedi_reimbursement_uploads
```

Jangan menjalankan `docker compose down -v` jika data masih dibutuhkan.

## Migration Gagal

Lihat log backend. Periksa tabel:

```cmd
docker exec -it jedi-reimbursement-mysql mysql ^
-uroot -proot jedi_reimbursement ^
-e "SELECT * FROM schema_migrations;"
```

Sesuaikan password jika sudah diubah.

## Smoke Test Gagal

Full Docker:

```cmd
smoke-test.bat
```

Mode pengembangan:

```cmd
powershell -NoProfile -ExecutionPolicy Bypass ^
-File smoke-test.ps1 ^
-FrontendUrl http://localhost:5173 ^
-ApiBaseUrl http://localhost:8080/api/v1
```

## Backup Kosong atau Gagal

Pastikan MySQL berjalan dan password di `.env.docker` benar.

```cmd
docker ps
```

## Restore Gagal

- Pastikan file SQL tersedia.
- Pastikan nama database benar.
- Pastikan container MySQL berjalan.
- Gunakan backup dari versi skema yang kompatibel.

## CSV Tidak Terbuka dengan Benar

Gunakan import UTF-8 pada Excel. Pemisah CSV mengikuti format standar koma dan
file menggunakan UTF-8 BOM.

## Frontend Menampilkan 404 Saat Refresh Route

Pada Full Docker, Nginx sudah memakai fallback SPA. Pastikan menggunakan
`frontend/nginx.conf` dari paket Part 8 dan image sudah dibangun ulang:

```cmd
docker-start.bat
```
