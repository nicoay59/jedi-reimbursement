# Frontend Professional 1.2.0

Frontend React menyediakan halaman karyawan, administrator, dashboard,
laporan, dan error fallback.

## Mode Lokal

```cmd
copy .env.example .env
npm install
npm run dev
```

## Build

```cmd
npm run build
```

## Docker

Dockerfile membangun React menggunakan:

```text
VITE_API_BASE_URL=/api/v1
```

Hasil build disajikan oleh Nginx. Route SPA memakai fallback ke
`index.html`, sedangkan `/api` diteruskan ke backend.
