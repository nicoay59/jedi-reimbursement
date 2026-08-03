@echo off
setlocal
cd /d "%~dp0"

if not exist "backend\.env" copy "backend\.env.example" "backend\.env" >nul

call "%~dp0start-database.bat"
if errorlevel 1 exit /b 1

cd /d "%~dp0backend"
go run ./cmd/migrate
if errorlevel 1 (
    echo [ERROR] Migration gagal.
    pause
    exit /b 1
)

go run ./cmd/seed
if errorlevel 1 (
    echo [ERROR] Seed gagal.
    pause
    exit /b 1
)

echo Backend berjalan di http://localhost:8080
go run ./cmd/api
pause
