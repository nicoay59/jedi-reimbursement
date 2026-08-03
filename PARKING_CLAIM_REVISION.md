# Revisi Klaim Parkir — v1.1.1

## Ketentuan yang Diterapkan

- Satu pengajuan hanya boleh mencakup rentang tanggal dalam satu bulan kalender.
- Nominal maksimal setiap pengajuan adalah Rp200.000.
- Batas Rp200.000 berlaku per pengajuan, bukan akumulasi seluruh klaim dalam bulan yang sama.
- Karyawan dapat membuat pengajuan terpisah untuk bulan yang berbeda.
- Klaim dapat diajukan untuk bulan berjalan sampai maksimal tiga bulan sebelumnya.
- Tanggal tidak boleh berada di masa depan.
- Beberapa bukti dalam satu pengajuan dapat digabung menjadi satu PDF maksimal 5 MB.

## Contoh

Jika tanggal saat ini 30 Juli 2026, karyawan dapat membuat pengajuan terpisah untuk:

```text
April 2026
Mei 2026
Juni 2026
Juli 2026
```

Contoh pengajuan yang valid:

```text
Pengajuan 1
Periode: 1–30 April 2026
Nominal: Rp200.000

Pengajuan 2
Periode: 1–31 Mei 2026
Nominal: Rp175.000

Pengajuan 3
Periode: 1–30 Juni 2026
Nominal: Rp200.000
```

Setiap pengajuan diproses sebagai klaim yang terpisah.

## Perubahan Teknis

- Menghapus perhitungan kuota akumulatif bulanan.
- Menghapus endpoint `/parking-claims/quota`.
- Mengubah validasi backend menjadi maksimal Rp200.000 per pengajuan.
- Menyederhanakan form frontend dan menghapus indikator pemakaian/sisa kuota.
- Migration `007_add_parking_claim_date_range.sql` tetap digunakan untuk menyimpan tanggal mulai dan tanggal selesai.
