@echo off
setlocal
cd /d "%~dp0"

docker compose --env-file .env.docker.example config >nul
if errorlevel 1 (
    echo [ERROR] Konfigurasi Docker Compose tidak valid.
    exit /b 1
)

call "%~dp0start-database.bat"
if errorlevel 1 exit /b 1

pushd backend
go mod tidy
gofmt -w cmd internal
go test ./...
if errorlevel 1 (
    popd
    echo [ERROR] Backend gagal diuji.
    exit /b 1
)
go build -o bin\jedi-reimbursement-api.exe .\cmd\api
if errorlevel 1 (
    popd
    echo [ERROR] Backend API gagal dibangun.
    exit /b 1
)

go build -o bin\jedi-reimbursement-migrate.exe .\cmd\migrate
if errorlevel 1 (
    popd
    echo [ERROR] Migration command gagal dibangun.
    exit /b 1
)

go build -o bin\jedi-reimbursement-seed.exe .\cmd\seed
if errorlevel 1 (
    popd
    echo [ERROR] Seed command gagal dibangun.
    exit /b 1
)

go build -o bin\jedi-reimbursement-wait-db.exe .\cmd\wait-db
if errorlevel 1 (
    popd
    echo [ERROR] Wait-db command gagal dibangun.
    exit /b 1
)

go run ./cmd/migrate
if errorlevel 1 (
    popd
    echo [ERROR] Migration gagal dijalankan.
    exit /b 1
)

go run ./cmd/seed
if errorlevel 1 (
    popd
    echo [ERROR] Seed gagal dijalankan.
    exit /b 1
)
popd

pushd frontend
call npm config set registry https://registry.npmjs.org/
if not exist "node_modules\.bin\vite.cmd" call npm install
call npm run build
if errorlevel 1 (
    popd
    echo [ERROR] Frontend gagal dibangun.
    exit /b 1
)
popd

echo [OK] Seluruh validasi Part 8 berhasil.
pause
