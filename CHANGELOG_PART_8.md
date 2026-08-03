# Changelog Part 8 — Finalisasi

## Docker dan Deployment

- Menambahkan Dockerfile multi-stage backend.
- Menambahkan Dockerfile multi-stage frontend.
- Menambahkan Nginx untuk SPA dan reverse proxy API.
- Menambahkan Docker Compose full stack.
- Menambahkan health check seluruh service.
- Menambahkan startup backend otomatis: tunggu database, migration, seed, API.
- Menambahkan Docker volume untuk MySQL dan bukti parkir.
- Menambahkan konfigurasi environment Docker.
- Menambahkan script start, stop, log, reset, backup, dan restore.

## Backend

- Menambahkan perintah `wait-db`.
- Menambahkan middleware security headers.
- Memperbarui versi health endpoint menjadi Part 8.
- Menambahkan pengujian security headers.
- Mempertahankan seluruh pengujian Part 1–7.

## Frontend

- Menambahkan React Error Boundary.
- Menambahkan halaman fallback kesalahan aplikasi.
- Build produksi memakai API relatif `/api/v1`.
- Nginx mendukung SPA route dan cache aset statis.

## Quality Assurance

- Menambahkan smoke test PowerShell dan batch wrapper.
- Menambahkan dokumentasi test case.
- Menambahkan release checklist.
- Menambahkan panduan screenshot dan narasi BAB IV.
- Menambahkan dokumentasi arsitektur, keamanan, deployment, dan manual pengguna.

## Versi Final

```text
1.0.0
```

## Penyempurnaan Antarmuka Produksi

- Mendesain ulang halaman login agar lebih profesional dan responsif.
- Menghapus label Part, badge Development, halaman status sistem, dan akun pengembangan dari UI.
- Menyempurnakan navigasi, profil pengguna, dan tombol keluar.
- Mendesain ulang dashboard karyawan dengan menu tindakan yang lebih jelas.
- Menyamakan warna, form, tombol, kartu, dan responsivitas seluruh aplikasi.
