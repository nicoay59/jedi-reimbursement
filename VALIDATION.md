# Validation Report — Corporate Professional UI 1.2.0

Tanggal validasi: 1 Agustus 2026

## Antarmuka

- Halaman login menggunakan tampilan enterprise perusahaan.
- Tidak ditemukan teks `Part`, `Development`, `Lihat status sistem`, atau tombol akun pengembangan pada source frontend.
- Navbar memiliki ikon dan indikator halaman aktif.
- Dashboard karyawan dan administrator telah disempurnakan.
- Form, tabel, badge, alert, timeline, empty state, dan loading state menggunakan gaya yang konsisten.
- Tampilan responsif tersedia untuk desktop, tablet, dan perangkat seluler.
- Dukungan pengurangan animasi tersedia melalui `prefers-reduced-motion`.

## Pemeriksaan source frontend

- Seluruh file JavaScript dan JSX berhasil diparse menggunakan TypeScript transpiler.
- File CSS berhasil diparse menggunakan PostCSS.
- Tidak ada perubahan endpoint atau aturan bisnis backend.

## Build frontend

`npm install` tidak dapat diselesaikan di lingkungan pembuatan paket karena
registry internal mengembalikan HTTP 404 untuk `@vitejs/plugin-react@6.0.3`.
Karena itu, build Vite penuh harus dijalankan pada komputer pengguna melalui:

```cmd
setup-windows.bat
```

atau:

```cmd
docker-start.bat
```

## Pengujian pada komputer pengguna

1. Jalankan aplikasi.
2. Periksa login pada desktop dan mobile.
3. Login sebagai karyawan dan buka semua menu klaim.
4. Login sebagai administrator dan buka dashboard, pemeriksaan, serta laporan.
5. Pastikan menu aktif mengikuti halaman yang dibuka.
6. Uji form, tabel, pagination, upload bukti, dan ekspor CSV.
