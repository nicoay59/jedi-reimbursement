# Release Checklist 1.0.0

## Source Code

- [ ] Tidak ada file `.go`, `.js`, atau `.jsx` kosong.
- [ ] `gofmt` sudah dijalankan.
- [ ] `go test ./...` berhasil.
- [ ] Backend berhasil dibangun.
- [ ] Frontend `npm run build` berhasil.
- [ ] Tidak ada `node_modules` di ZIP.
- [ ] Tidak ada `.env` rahasia di ZIP.

## Environment

- [ ] `JWT_SECRET` sudah diganti.
- [ ] Password database sudah diganti.
- [ ] Password admin sudah diganti.
- [ ] Password karyawan demo sudah diganti atau dihapus.
- [ ] `FRONTEND_URL` sesuai domain.
- [ ] Time zone sesuai.

## Docker

- [ ] `docker compose config` berhasil.
- [ ] MySQL healthy.
- [ ] Backend healthy.
- [ ] Frontend healthy.
- [ ] Migration berhasil.
- [ ] Seed berhasil.
- [ ] Volume database aktif.
- [ ] Volume upload aktif.

## Functional Test

- [ ] Login admin.
- [ ] Login karyawan.
- [ ] Buat klaim parkir.
- [ ] Unduh bukti.
- [ ] Buat klaim lembur.
- [ ] Setujui klaim.
- [ ] Tolak klaim.
- [ ] Riwayat status tampil.
- [ ] Dashboard tampil.
- [ ] Laporan tampil.
- [ ] CSV berhasil dibuka.

## Operational Test

- [ ] `smoke-test.bat` berhasil.
- [ ] Backup database berhasil.
- [ ] Restore diuji.
- [ ] Restart container tidak menghapus data.
- [ ] Log tidak menunjukkan error berulang.

## Dokumentasi

- [ ] README diperiksa.
- [ ] Manual pengguna diperiksa.
- [ ] API diperiksa.
- [ ] Test case dilengkapi hasil aktual.
- [ ] Screenshot BAB IV dibuat.
- [ ] Versi release dicatat.
