# Panduan Deployment

## Mode Pengembangan Windows

Prasyarat:

- Go 1.23 atau lebih baru.
- Node.js sesuai `frontend/package.json`.
- Docker Desktop.
- Port 3306, 8080, dan 5173 tersedia.

Jalankan:

```cmd
setup-windows.bat
run-all.bat
```

Akses:

```text
Frontend : http://localhost:5173
Backend  : http://localhost:8080
MySQL    : localhost:3306
```

## Mode Full Docker

Prasyarat:

- Docker Desktop dengan Docker Compose v2.
- Port 8088 dan 3306 tersedia.

### 1. Siapkan Environment

Pada eksekusi pertama, `docker-start.bat` membuat:

```text
.env.docker
```

Buka file tersebut dan ganti minimal:

```text
MYSQL_ROOT_PASSWORD
DB_PASSWORD
JWT_SECRET
ADMIN_PASSWORD
EMPLOYEE_PASSWORD
```

### 2. Jalankan

```cmd
docker-start.bat
```

Akses:

```text
http://localhost:8088
```

### 3. Lihat Log

```cmd
docker-logs.bat
```

### 4. Hentikan

```cmd
docker-stop.bat
```

Perintah stop tidak menghapus volume.

### 5. Reset Total

```cmd
docker-reset.bat
```

Reset menghapus database dan bukti upload. Script meminta konfirmasi `HAPUS`.

## Proses Startup Backend Docker

Container backend menjalankan urutan:

```text
wait-db
→ migrate
→ seed
→ api
```

Jika migration gagal, API tidak dijalankan. Hal ini mencegah aplikasi berjalan
dengan skema database yang belum sesuai.

## Health Check

MySQL:

```text
mysqladmin ping
```

Backend:

```text
GET /api/v1/health
```

Frontend:

```text
GET /healthz
```

## Deployment Server

Untuk server produksi:

1. Gunakan domain.
2. Pasang TLS/HTTPS di reverse proxy.
3. Ubah `FRONTEND_URL` sesuai domain.
4. Jangan expose port MySQL ke publik.
5. Gunakan password kuat.
6. Gunakan Docker volume atau storage persisten.
7. Jadwalkan backup.
8. Batasi firewall hanya pada port HTTP/HTTPS.
9. Pantau log dan kapasitas disk.
10. Uji restore backup secara berkala.

## Mengubah Port Frontend Docker

Pada `.env.docker`:

```text
FRONTEND_PORT=8090
```

Kemudian:

```cmd
docker-stop.bat
docker-start.bat
```

Akses menjadi:

```text
http://localhost:8090
```

## Backup

```cmd
backup-database.bat
```

## Restore

```cmd
restore-database.bat backups\nama-file.sql
```

Lakukan backup sebelum upgrade atau reset.
