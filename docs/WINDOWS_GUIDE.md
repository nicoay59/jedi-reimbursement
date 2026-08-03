# Panduan Windows Final

## Mode Pengembangan

Prasyarat:

```text
Go
Node.js
npm
Docker Desktop
```

Jalankan:

```cmd
setup-windows.bat
run-all.bat
```

Alamat:

```text
http://localhost:5173
```

Tutup jendela backend dan frontend untuk menghentikan mode pengembangan.

## Mode Full Docker

Jalankan:

```cmd
docker-start.bat
```

Alamat default:

```text
http://localhost:8088
```

Script akan:

1. Membuat `.env.docker` jika belum ada.
2. Mendeteksi volume MySQL dari container lama.
3. Menggunakan kembali volume jika tersedia.
4. Membangun backend dan frontend.
5. Menjalankan migration dan seed.
6. Menunggu seluruh service sehat.

Hentikan:

```cmd
docker-stop.bat
```

Lihat log:

```cmd
docker-logs.bat
```

## Smoke Test

```cmd
smoke-test.bat
```

## Backup

```cmd
backup-database.bat
```

## Restore

```cmd
restore-database.bat backups\nama-file.sql
```

## Reset

```cmd
docker-reset.bat
```

Reset menghapus database dan bukti pada Docker volume.

## Konflik Port

Part 8 tidak mempunyai `check-ports.bat`.

Periksa manual:

```cmd
netstat -ano | findstr :3306
netstat -ano | findstr :8080
netstat -ano | findstr :5173
netstat -ano | findstr :8088
```

Hentikan PID yang benar:

```cmd
taskkill /PID NOMOR_PID /F
```
