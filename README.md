# Reimbursement Portal — Corporate Professional Build

Part 8 adalah paket final kumulatif dari Part 1 sampai Part 7. Proyek sudah
mencakup autentikasi, klaim parkir, klaim lembur, persetujuan administrator,
dashboard, laporan, ekspor CSV, dokumentasi Kerja Praktik, serta dua cara
menjalankan aplikasi: mode pengembangan dan mode Docker penuh.

## Versi

```text
1.2.0
```


## Corporate Professional UI

Versi 1.2.0 menggunakan tampilan perusahaan yang lebih formal dengan corporate
navigation, active menu, dashboard berbasis ikon, form dan tabel yang konsisten,
serta desain responsif untuk penggunaan internal PT. Jedi Global Teknologi.

## Teknologi

- Backend: Golang 1.23+
- Frontend: React.js 19
- Build tool: Vite 8
- UI: Bootstrap 5.3
- Database: MySQL 8
- Web server produksi: Nginx Alpine
- Autentikasi: JWT HS256
- Container: Docker Compose

## Modul Final

### Karyawan

- Login dan logout.
- Pengajuan klaim parkir.
- Unggah bukti JPG, PNG, atau PDF.
- Pengajuan klaim lembur.
- Perhitungan durasi otomatis.
- Riwayat dan detail klaim.
- Status dan catatan administrator.

### Administrator

- Dashboard statistik berdasarkan periode.
- Pemeriksaan klaim parkir dan lembur.
- Persetujuan dan penolakan.
- Catatan administrator.
- Riwayat perubahan status.
- Laporan dengan filter dan pagination.
- Ekspor CSV untuk Microsoft Excel.

### Operasional

- Migration dan seed otomatis.
- Full-stack Docker Compose.
- Health check MySQL, backend, dan frontend.
- Reverse proxy `/api` melalui Nginx.
- Penyimpanan bukti menggunakan Docker volume.
- Backup dan restore database.
- Smoke test otomatis.
- Security headers.
- Dokumentasi deployment, pengguna, pengujian, dan BAB IV.

## Kebijakan Klaim Parkir per Pengajuan

- Klaim menggunakan rentang tanggal dalam satu bulan kalender.
- Nominal maksimal Rp200.000 untuk setiap pengajuan.
- Pengajuan untuk bulan berbeda dibuat secara terpisah, dengan periode bulan berjalan sampai tiga bulan sebelumnya.
- Bukti transaksi dapat digabung menjadi satu PDF maksimal 5 MB.

## Akun Pengembangan

Administrator:

```text
Email    : admin@jedi.local
Password : Admin123!
```

Karyawan:

```text
Email    : employee@jedi.local
Password : Employee123!
```

Ganti seluruh password dan `JWT_SECRET` sebelum deployment nyata.

# Pilihan Menjalankan

## Mode A — Pengembangan Windows

Mode ini menjalankan:

- MySQL melalui Docker.
- Backend langsung dengan Go.
- Frontend langsung dengan Node.js dan Vite.

Tutup terminal Part 7, buka Docker Desktop, kemudian:

```cmd
setup-windows.bat
run-all.bat
```

Buka:

```text
http://localhost:5173
```

Mode ini tetap tidak menggunakan `check-ports.bat`.

## Mode B — Full Docker

Mode ini menjalankan MySQL, backend, migration, seed, frontend, Nginx, health
check, dan penyimpanan bukti melalui Docker Compose.

Jalankan:

```cmd
docker-start.bat
```

Buka:

```text
http://localhost:8088
```

Hentikan:

```cmd
docker-stop.bat
```

# Struktur Docker

```text
Browser
   |
   v
Nginx + React :8088
   |
   +---- /api/* ----> Golang API :8080
                         |
                         v
                      MySQL :3306
```

Backend tidak dipublikasikan langsung pada mode Docker penuh. Browser mengakses
API melalui reverse proxy Nginx.

# Pemeriksaan Akhir

Setelah Full Docker aktif:

```cmd
smoke-test.bat
```

Smoke test memeriksa:

- Frontend.
- Health API.
- Login administrator.
- Endpoint profil.
- Dashboard administrator.
- Laporan administrator.

# Backup Database

```cmd
backup-database.bat
```

Hasil backup disimpan pada folder:

```text
backups
```

Restore:

```cmd
restore-database.bat backups\nama-file.sql
```

# Dokumentasi

Lihat folder `docs`:

- `USER_MANUAL.md`
- `DEPLOYMENT.md`
- `ARCHITECTURE.md`
- `SECURITY.md`
- `TEST_CASES.md`
- `BAB_IV_GUIDE.md`
- `API.md`
- `RELEASE_CHECKLIST.md`
- `TROUBLESHOOTING.md`

## Catatan Final

Part 8 tidak menambahkan tabel baru. Seluruh migration Part 1–7 tetap digunakan.
Paket final dirancang untuk bahan Kerja Praktik, demonstrasi, pengujian, dan
pengembangan berikutnya.
