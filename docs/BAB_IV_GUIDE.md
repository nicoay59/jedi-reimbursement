# Panduan Penyusunan BAB IV

Dokumen ini dapat digunakan sebagai kerangka penulisan hasil dan pembahasan
implementasi sistem.

## 4.1 Hasil Analisis Kebutuhan

Jelaskan dua aktor:

### Karyawan

- Login.
- Mengajukan klaim parkir.
- Mengunggah bukti pembayaran.
- Mengajukan klaim lembur.
- Melihat riwayat dan status.

### Administrator

- Login.
- Memeriksa pengajuan.
- Menyetujui atau menolak.
- Melihat dashboard.
- Membuka laporan.
- Mengekspor CSV.

Screenshot yang disarankan:

1. Halaman login.
2. Dashboard karyawan.
3. Dashboard administrator.

## 4.2 Implementasi Basis Data

Bahas tabel:

- `users`
- `parking_claims`
- `overtime_claims`
- `claim_status_histories`
- `schema_migrations`

Screenshot:

1. Daftar tabel MySQL.
2. Struktur tabel `users`.
3. Struktur `parking_claims`.
4. Struktur `overtime_claims`.
5. Data riwayat status.

Contoh narasi:

> Basis data MySQL digunakan untuk menyimpan data pengguna, pengajuan klaim
> parkir, pengajuan lembur, dan riwayat perubahan status. Migration digunakan
> agar perubahan struktur basis data tercatat dan dapat diterapkan secara
> konsisten.

## 4.3 Implementasi Backend

Bahas struktur:

```text
handler → service → repository → database
```

Fitur yang dibahas:

- JWT.
- Middleware role.
- Validasi request.
- Upload bukti.
- Perhitungan durasi.
- Transaksi approval.
- Laporan.
- Ekspor CSV.

Screenshot:

1. Struktur folder backend.
2. Terminal API berjalan.
3. Respons health endpoint.
4. Respons login.
5. Respons pembuatan klaim.
6. Respons approval.

## 4.4 Implementasi Frontend

Bahas:

- React component.
- AuthContext.
- Protected route.
- Role route.
- Form.
- Dashboard.
- Laporan.
- Error Boundary.

Screenshot aplikasi:

1. Login.
2. Form parkir.
3. Detail parkir.
4. Form lembur.
5. Riwayat lembur.
6. Daftar pemeriksaan.
7. Detail approval.
8. Dashboard statistik.
9. Laporan.
10. CSV di Excel.

## 4.5 Implementasi Docker

Bahas service:

- MySQL.
- Backend.
- Frontend/Nginx.

Screenshot:

1. `docker compose ps`.
2. Docker Desktop menampilkan tiga container.
3. Health status `healthy`.
4. Aplikasi pada `localhost:8088`.
5. Output `smoke-test.bat`.

Contoh narasi:

> Docker Compose digunakan untuk menyatukan proses deployment database,
> backend, dan frontend. Backend menjalankan proses menunggu database,
> migration, seed, kemudian API secara otomatis.

## 4.6 Pengujian Sistem

Gunakan tabel pada `TEST_CASES.md`.

Kelompok pengujian:

- Autentikasi.
- Klaim parkir.
- Klaim lembur.
- Persetujuan.
- Laporan.
- Docker.
- Backup dan restore.

## 4.7 Pembahasan

Bahas hasil:

- Pemisahan role berjalan.
- Validasi mengurangi data salah.
- Transaksi menjaga konsistensi approval.
- Dashboard mempermudah pemantauan.
- CSV membantu rekap.
- Docker mempermudah instalasi.

## 4.8 Keterbatasan

Tuliskan dengan jujur:

- Belum ada email notifikasi.
- Belum ada refresh token.
- Belum ada approval bertingkat.
- Belum ada object storage.
- Belum ada manajemen pengguna melalui UI.

## Checklist Screenshot

Gunakan nama file terurut:

```text
01-login.png
02-dashboard-karyawan.png
03-form-parkir.png
04-detail-parkir.png
05-form-lembur.png
06-riwayat-lembur.png
07-daftar-pemeriksaan.png
08-detail-approval.png
09-dashboard-admin.png
10-laporan.png
11-csv-excel.png
12-docker-compose.png
13-smoke-test.png
```

Pastikan tidak menampilkan password, JWT, atau secret pada screenshot.
