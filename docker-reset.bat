@echo off
setlocal
cd /d "%~dp0"

echo PERINGATAN:
echo Perintah ini menghapus container, database, dan bukti upload Docker.
echo Semua data pada Docker volume akan hilang.
echo.
set /p CONFIRM=Ketik HAPUS untuk melanjutkan: 

if /I not "%CONFIRM%"=="HAPUS" (
    echo Reset dibatalkan.
    pause
    exit /b 0
)

if exist ".env.docker" (
    docker compose --env-file .env.docker down -v --remove-orphans
) else (
    docker compose down -v --remove-orphans
)

echo Seluruh data Docker proyek telah dihapus.
pause
