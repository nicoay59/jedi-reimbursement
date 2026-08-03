@echo off
setlocal
cd /d "%~dp0"

echo ============================================
echo Setup Jedi Reimbursement - Part 8 Final Clean
echo ============================================

for %%C in (go node npm docker) do (
    where %%C >nul 2>nul
    if errorlevel 1 (
        echo [ERROR] %%C tidak ditemukan pada PATH.
        exit /b 1
    )
)

docker info >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Docker Desktop belum berjalan.
    exit /b 1
)

if not exist "backend\.env" copy "backend\.env.example" "backend\.env" >nul
if not exist "frontend\.env" copy "frontend\.env.example" "frontend\.env" >nul

if not exist "backend\storage\uploads\parking-receipts" mkdir "backend\storage\uploads\parking-receipts"

call "%~dp0start-database.bat"
if errorlevel 1 exit /b 1

echo.
echo Menyiapkan backend...
pushd backend
go mod tidy
if errorlevel 1 (
    popd
    echo [ERROR] go mod tidy gagal.
    exit /b 1
)

gofmt -w cmd internal
go test ./...
if errorlevel 1 (
    popd
    echo [ERROR] Pengujian backend gagal.
    exit /b 1
)

if not exist "bin" mkdir bin
go build -o bin\jedi-reimbursement-api.exe .\cmd\api
if errorlevel 1 (
    popd
    echo [ERROR] Build backend gagal.
    exit /b 1
)

go build -o bin\jedi-reimbursement-wait-db.exe .\cmd\wait-db
if errorlevel 1 (
    popd
    echo [ERROR] Build wait-db gagal.
    exit /b 1
)

go run ./cmd/migrate
if errorlevel 1 (
    popd
    echo [ERROR] Migration gagal.
    exit /b 1
)

go run ./cmd/seed
if errorlevel 1 (
    popd
    echo [ERROR] Seed gagal.
    exit /b 1
)
popd

echo.
echo Menyiapkan frontend...
pushd frontend
if exist "package-lock.json" (
    findstr /i /c:"internal.api.openai.org" package-lock.json >nul
    if not errorlevel 1 del /f /q package-lock.json
)

call npm config set registry https://registry.npmjs.org/
if exist "node_modules" rmdir /s /q node_modules
call npm install
if errorlevel 1 (
    popd
    echo [ERROR] npm install gagal.
    exit /b 1
)

call npm run build
if errorlevel 1 (
    popd
    echo [ERROR] Build frontend gagal.
    exit /b 1
)
popd

echo.
echo ============================================
echo Setup Part 8 berhasil.
echo Jalankan run-all.bat
echo ============================================
pause
