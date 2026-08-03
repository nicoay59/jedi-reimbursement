# API Final — Versi 1.2.0

## Base URL

Mode pengembangan:

```text
http://localhost:8080/api/v1
```

Mode Full Docker:

```text
http://localhost:8088/api/v1
```

## Format Respons

```json
{
  "success": true,
  "message": "Pesan",
  "data": {}
}
```

Respons validasi dapat memiliki:

```json
{
  "success": false,
  "message": "Data tidak valid",
  "errors": {
    "field": "Pesan validasi"
  }
}
```

## Health

```http
GET /health
```

## Autentikasi

### Login

```http
POST /auth/login
Content-Type: application/json
```

```json
{
  "email": "admin@jedi.local",
  "password": "Admin123!"
}
```

### Profil

```http
GET /auth/me
Authorization: Bearer <token>
```

### Logout

```http
POST /auth/logout
Authorization: Bearer <token>
```

## Klaim Parkir — EMPLOYEE

```text
POST /parking-claims
GET  /parking-claims?page=1&limit=10
GET  /parking-claims/{id}
GET  /parking-claims/{id}/receipt
```

Request create menggunakan `multipart/form-data`:

```text
parking_start_date
parking_end_date
parking_location
amount
description
receipt
```

Aturan klaim parkir:

- rentang tanggal harus berada pada bulan kalender yang sama;
- tanggal tidak boleh di masa depan;
- klaim hanya dapat diajukan untuk bulan berjalan sampai tiga bulan sebelumnya;
- nominal maksimal setiap pengajuan adalah Rp200.000;
- pengajuan bulan yang berbeda dibuat secara terpisah;
- tidak ada perhitungan akumulasi kuota bulanan antar-pengajuan.


## Klaim Lembur — EMPLOYEE

```text
POST /overtime-claims
GET  /overtime-claims?page=1&limit=10
GET  /overtime-claims/{id}
```

Request create:

```json
{
  "overtime_date": "2026-07-18",
  "start_time": "18:00",
  "end_time": "21:30",
  "work_description": "Menyelesaikan laporan bulanan"
}
```

## Pemeriksaan — ADMIN

```text
GET   /admin/claims
GET   /admin/claims/{type}/{id}
PATCH /admin/claims/{type}/{id}/status
GET   /admin/claims/{type}/{id}/history
GET   /admin/claims/{type}/{id}/receipt
```

Filter daftar:

```text
type=ALL|PARKING|OVERTIME
status=ALL|PENDING|APPROVED|REJECTED
page=1
limit=10
```

Request keputusan:

```json
{
  "status": "APPROVED",
  "note": "Data telah sesuai"
}
```

atau:

```json
{
  "status": "REJECTED",
  "note": "Bukti tidak dapat dibaca"
}
```

## Dashboard — ADMIN

```http
GET /admin/dashboard?start_date=2026-07-01&end_date=2026-07-18
```

## Laporan — ADMIN

```http
GET /admin/reports
    ?start_date=2026-07-01
    &end_date=2026-07-18
    &type=ALL
    &status=ALL
    &page=1
    &limit=10
```

## Ekspor — ADMIN

```http
GET /admin/reports/export
    ?start_date=2026-07-01
    &end_date=2026-07-18
    &type=ALL
    &status=ALL
```

## Status HTTP

| Status | Arti |
|---|---|
| 200 | Berhasil |
| 201 | Data berhasil dibuat |
| 204 | Tidak ada isi |
| 400 | Request tidak valid |
| 401 | Belum login atau token salah |
| 403 | Role tidak diizinkan |
| 404 | Data tidak ditemukan |
| 409 | Klaim sudah diputuskan |
| 413 | File terlalu besar |
| 422 | Validasi gagal |
| 500 | Kesalahan server |
