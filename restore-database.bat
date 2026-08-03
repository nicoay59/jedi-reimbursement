@echo off
setlocal EnableExtensions
cd /d "%~dp0"

if "%~1"=="" (
    echo Penggunaan:
    echo restore-database.bat backups\nama-file.sql
    pause
    exit /b 1
)

if not exist "%~1" (
    echo [ERROR] File backup tidak ditemukan:
    echo %~1
    pause
    exit /b 1
)

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

echo PERINGATAN:
echo Database %DB_NAME% akan ditimpa menggunakan:
echo %~1
echo.
set /p CONFIRM=Ketik RESTORE untuk melanjutkan: 

if /I not "%CONFIRM%"=="RESTORE" (
    echo Restore dibatalkan.
    pause
    exit /b 0
)

type "%~1" | docker exec -i jedi-reimbursement-mysql sh -c "exec mysql -uroot -p\"%DB_ROOT_PASSWORD%\" %DB_NAME%"

if errorlevel 1 (
    echo [ERROR] Restore database gagal.
    pause
    exit /b 1
)

echo [OK] Restore database berhasil.
pause
