# Test Cases Final

## Autentikasi

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| AUTH-01 | Login admin valid | Masuk dashboard admin |
| AUTH-02 | Login karyawan valid | Masuk dashboard karyawan |
| AUTH-03 | Password salah | HTTP 401 dan pesan kesalahan |
| AUTH-04 | Token tidak ada | Endpoint privat menolak |
| AUTH-05 | Role salah | HTTP 403 |

## Klaim Parkir

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| PK-01 | Form valid + JPG | Klaim tersimpan |
| PK-02 | PDF valid | Klaim tersimpan |
| PK-03 | File lebih dari 5 MB | Ditolak |
| PK-04 | File TXT | Ditolak |
| PK-05 | Tanggal masa depan | Ditolak |
| PK-06 | Buka klaim pengguna lain | Tidak ditemukan |
| PK-07 | Unduh bukti | File berhasil diunduh |

## Klaim Lembur

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| OT-01 | 18:00–21:30 | Durasi 3,5 jam |
| OT-02 | 22:00–01:30 | Durasi 3,5 jam |
| OT-03 | Durasi kurang 30 menit | Ditolak |
| OT-04 | Durasi lebih 16 jam | Ditolak |
| OT-05 | Deskripsi terlalu pendek | Ditolak |

## Persetujuan

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| RV-01 | Setujui PENDING | Status APPROVED |
| RV-02 | Tolak dengan catatan | Status REJECTED |
| RV-03 | Tolak tanpa catatan | Ditolak validasi |
| RV-04 | Proses ulang APPROVED | HTTP 409 |
| RV-05 | Riwayat status | Riwayat tampil berurutan |

## Dashboard dan Laporan

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| RP-01 | Dashboard bulan berjalan | Statistik tampil |
| RP-02 | Filter periode | Data mengikuti periode |
| RP-03 | Filter jenis PARKING | Hanya parkir |
| RP-04 | Filter status APPROVED | Hanya disetujui |
| RP-05 | Ekspor CSV | File CSV terunduh |
| RP-06 | Rentang lebih 366 hari | Ditolak |

## Docker

| ID | Pengujian | Hasil yang Diharapkan |
|---|---|---|
| DK-01 | `docker-start.bat` | Semua service healthy |
| DK-02 | Restart container | Data tetap tersedia |
| DK-03 | Backup database | File SQL dibuat |
| DK-04 | Restore backup | Data kembali |
| DK-05 | `smoke-test.bat` | Semua pemeriksaan berhasil |

## Kriteria Kelulusan

- Semua test kritis berhasil.
- Tidak ada error build.
- Tidak ada file source kosong.
- Migration dapat dijalankan dari database kosong.
- Data tetap tersedia setelah restart.
- Backup dan restore dapat digunakan.
