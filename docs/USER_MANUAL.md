# Manual Pengguna

## Login

Buka aplikasi lalu masukkan email dan password.

Akun karyawan diarahkan ke dashboard karyawan. Akun administrator diarahkan ke
dashboard administrator.

## Karyawan

### Mengajukan Klaim Parkir

1. Pilih **Ajukan → Klaim Parkir**.
2. Pilih tanggal mulai dan tanggal selesai pada bulan yang sama.
3. Isi lokasi atau area parkir.
4. Isi total nominal periode tersebut.
5. Isi rincian klaim.
6. Unggah satu bukti JPG/PNG atau satu PDF gabungan.
7. Periksa kembali data, lalu klik **Kirim pengajuan**.

Ketentuan:

- satu pengajuan maksimal mencakup satu bulan kalender;
- nominal maksimal Rp200.000 untuk setiap pengajuan;
- periode maksimal bulan berjalan sampai tiga bulan sebelumnya;
- pengajuan untuk bulan yang berbeda dibuat secara terpisah;
- batas Rp200.000 berlaku pada setiap pengajuan, bukan total seluruh pengajuan dalam bulan tersebut.

Status awal:

```text
Menunggu
```

### Melihat Riwayat Parkir

1. Pilih menu **Parkir**.
2. Gunakan pagination jika data banyak.
3. Klik **Detail**.
4. Lihat status, catatan admin, dan bukti.

### Mengajukan Klaim Lembur

1. Pilih **Ajukan → Klaim Lembur**.
2. Isi tanggal.
3. Isi waktu mulai.
4. Isi waktu selesai.
5. Periksa estimasi durasi.
6. Isi deskripsi pekerjaan.
7. Klik **Kirim pengajuan**.

Jika waktu selesai lebih kecil dari waktu mulai, sistem menganggap lembur
berakhir pada hari berikutnya.

### Melihat Riwayat Lembur

1. Pilih menu **Lembur**.
2. Klik **Detail** pada pengajuan.
3. Lihat durasi, status, dan catatan administrator.

## Administrator

### Dashboard

Dashboard menampilkan:

- Total klaim.
- Klaim menunggu.
- Klaim disetujui.
- Klaim ditolak.
- Nominal parkir.
- Durasi lembur.
- Tren harian.
- Klaim terbaru.

Ubah tanggal mulai dan selesai lalu klik **Terapkan**.

### Pemeriksaan Klaim

1. Pilih **Pemeriksaan Klaim**.
2. Gunakan filter jenis dan status.
3. Klik **Periksa**.
4. Periksa data karyawan dan pengajuan.
5. Unduh bukti jika klaim parkir.
6. Pilih setujui atau tolak.
7. Isi catatan.
8. Klik tombol keputusan.

Catatan wajib saat menolak.

### Laporan

1. Pilih **Laporan**.
2. Pilih periode.
3. Pilih jenis.
4. Pilih status.
5. Klik **Terapkan**.
6. Klik **Ekspor CSV** untuk mengunduh laporan.

## Logout

Klik tombol **Logout** pada kanan atas. Token lokal dihapus dan pengguna kembali
ke halaman login.
