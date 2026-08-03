@echo off
setlocal EnableExtensions
cd /d "%~dp0"

docker ps --format "{{.Names}}" | findstr /x "jedi-reimbursement-mysql" >nul
if errorlevel 1 (
    echo [ERROR] Container jedi-reimbursement-mysql tidak berjalan.
    pause
    exit /b 1
)

set DB_ROOT_PASSWORD=root
if exist ".env.docker" (
    for /f "usebackq tokens=1,* delims==" %%A in (".env.docker") do (
        if /I "%%A"=="MYSQL_ROOT_PASSWORD" set DB_ROOT_PASSWORD=%%B
    )
)

set DB_NAME=jedi_reimbursement
if exist ".env.docker" (
    for /f "usebackq tokens=1,* delims==" %%A in (".env.docker") do (
        if /I "%%A"=="DB_NAME" set DB_NAME=%%B
    )
)

if not exist "backups" mkdir backups

for /f %%T in ('powershell -NoProfile -Command "Get-Date -Format yyyyMMdd-HHmmss"') do set TIMESTAMP=%%T
set BACKUP_FILE=backups\jedi-reimbursement-%TIMESTAMP%.sql

echo Membuat backup %BACKUP_FILE%...
docker exec jedi-reimbursement-mysql sh -c "exec mysqldump -uroot -p\"%DB_ROOT_PASSWORD%\" --single-transaction --routines --triggers %DB_NAME%" > "%BACKUP_FILE%"

if errorlevel 1 (
    if exist "%BACKUP_FILE%" del /f /q "%BACKUP_FILE%"
    echo [ERROR] Backup database gagal.
    pause
    exit /b 1
)

echo [OK] Backup berhasil:
echo %BACKUP_FILE%
pause
