# Arsitektur Sistem

## Gambaran Umum

Jedi Reimbursement menggunakan arsitektur client-server dengan REST API.

```text
React SPA
   |
   | HTTPS/HTTP + JSON
   v
Nginx Reverse Proxy
   |
   | /api/v1
   v
Golang REST API
   |
   | database/sql
   v
MySQL 8
```

## Backend

Struktur backend menggunakan pemisahan tanggung jawab:

```text
cmd
├── api
├── migrate
├── seed
└── wait-db

internal
├── config
├── database
├── dto
├── handlers
├── middleware
├── models
├── repositories
├── responses
├── routes
├── security
├── services
└── storage
```

### Alur Request

```text
HTTP Request
→ Recovery
→ Request ID
→ Logger
→ Security Headers
→ CORS
→ Authentication
→ Role Authorization
→ Handler
→ Service
→ Repository
→ MySQL
```

### Lapisan

- **Handler:** menerima dan mengembalikan HTTP.
- **Service:** menjalankan aturan bisnis dan validasi.
- **Repository:** menjalankan query database.
- **DTO:** menentukan format request dan response.
- **Model:** mewakili data domain.
- **Middleware:** autentikasi, role, logging, CORS, dan keamanan.
- **Storage:** menyimpan serta membaca bukti parkir.

## Frontend

Frontend merupakan React Single Page Application.

```text
src
├── components
├── constants
├── context
├── hooks
├── pages
├── routes
├── services
├── styles
└── utils
```

- `AuthContext` menyimpan sesi pengguna.
- `ProtectedRoute` membatasi pengguna yang belum login.
- `RoleRoute` membatasi halaman berdasarkan role.
- Service frontend mengakses REST API.
- Error Boundary menangani error render yang tidak terduga.

## Database

Tabel utama:

- `users`
- `parking_claims`
- `overtime_claims`
- `claim_status_histories`
- `schema_migrations`

Relasi utama:

```text
users 1 ----- n parking_claims
users 1 ----- n overtime_claims
users 1 ----- n claim_status_histories
```

## Deployment Docker

```text
jedi-reimbursement-frontend
├── Nginx
└── React static build

jedi-reimbursement-backend
├── wait-db
├── migrate
├── seed
└── API

jedi-reimbursement-mysql
└── MySQL 8
```

Volume:

- `jedi_reimbursement_mysql_data`
- `jedi_reimbursement_uploads`
