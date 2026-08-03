# Keamanan Sistem

## Autentikasi

- Password disimpan menggunakan PBKDF2-HMAC-SHA256.
- Login menghasilkan JWT HS256.
- Endpoint privat membutuhkan Bearer token.
- JWT memiliki waktu kedaluwarsa.
- Akun tidak aktif ditolak oleh backend.

## Otorisasi

- Role `EMPLOYEE` mengakses klaim miliknya.
- Role `ADMIN` mengakses pemeriksaan, dashboard, dan laporan.
- Karyawan tidak dapat membuka klaim pengguna lain.
- Klaim yang sudah diputuskan tidak dapat diproses ulang.

## Upload Bukti

- Format hanya JPG, PNG, atau PDF.
- MIME dideteksi berdasarkan isi file.
- Ukuran file dibatasi.
- Nama file dibuat secara acak.
- Path traversal ditolak.
- Folder upload tidak disajikan sebagai folder publik.
- Bukti dibaca melalui endpoint terautentikasi.

## CSV

- CSV menggunakan escaping standar.
- Sel yang dimulai `=`, `+`, `-`, atau `@` diberi apostrof.
- Hal ini mengurangi risiko formula injection saat dibuka di spreadsheet.

## HTTP Headers

Backend mengirim:

```text
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
Referrer-Policy: no-referrer
Permissions-Policy: camera=(), microphone=(), geolocation=()
Cache-Control: no-store
```

Nginx menambahkan Content Security Policy dan header keamanan untuk frontend.

## Rekomendasi Produksi

Sebelum deployment nyata:

1. Ganti `JWT_SECRET`.
2. Ganti seluruh password seed.
3. Gunakan HTTPS.
4. Jangan publikasikan port MySQL ke internet.
5. Batasi akses server dan volume.
6. Lakukan backup terjadwal.
7. Simpan secret pada secret manager.
8. Tambahkan rate limiting pada login.
9. Pertimbangkan refresh token dan pencabutan token.
10. Gunakan object storage untuk bukti jika skala meningkat.
