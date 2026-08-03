@echo off
setlocal
cd /d "%~dp0"

where docker >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Docker tidak ditemukan.
    pause
    exit /b 1
)

docker info >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Docker Desktop belum berjalan.
    pause
    exit /b 1
)

docker ps -a --format "{{.Names}}" | findstr /x "jedi-reimbursement-mysql" >nul
if not errorlevel 1 (
    echo Container MySQL sudah tersedia. Menjalankan container...
    docker start jedi-reimbursement-mysql >nul
) else (
    echo Membuat dan menjalankan container MySQL...
    docker compose up -d mysql
)

if errorlevel 1 (
    echo [ERROR] MySQL gagal dijalankan.
    pause
    exit /b 1
)

echo Menunggu MySQL siap...
set /a ATTEMPT=0

:WAIT_MYSQL
set /a ATTEMPT+=1
docker exec jedi-reimbursement-mysql mysqladmin ping -h localhost -uroot -proot --silent >nul 2>nul
if not errorlevel 1 goto MYSQL_READY

if %ATTEMPT% GEQ 40 (
    echo [ERROR] MySQL belum siap setelah sekitar 2 menit.
    echo Periksa: docker logs jedi-reimbursement-mysql
    pause
    exit /b 1
)

timeout /t 3 /nobreak >nul
goto WAIT_MYSQL

:MYSQL_READY
echo [OK] MySQL siap pada localhost:3306.
