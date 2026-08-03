@echo off
setlocal EnableExtensions EnableDelayedExpansion
cd /d "%~dp0"

where docker >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Docker tidak ditemukan pada PATH.
    pause
    exit /b 1
)

docker info >nul 2>nul
if errorlevel 1 (
    echo [ERROR] Docker Desktop belum berjalan.
    pause
    exit /b 1
)

if not exist ".env.docker" (
    copy ".env.docker.example" ".env.docker" >nul
    echo [INFO] .env.docker dibuat dari contoh.
    echo [INFO] Ubah password dan JWT_SECRET sebelum deployment nyata.
)

set FRONTEND_PORT=8088
for /f "usebackq tokens=1,* delims==" %%A in (".env.docker") do (
    if /I "%%A"=="FRONTEND_PORT" set FRONTEND_PORT=%%B
)

set EXISTING_MYSQL_VOLUME=
docker ps -a --format "{{.Names}}" | findstr /x "jedi-reimbursement-mysql" >nul
if not errorlevel 1 (
    for /f "usebackq delims=" %%V in (`docker inspect --format="{{(index .Mounts 0).Name}}" jedi-reimbursement-mysql 2^>nul`) do set EXISTING_MYSQL_VOLUME=%%V

    if defined EXISTING_MYSQL_VOLUME (
        echo [INFO] Menggunakan kembali volume MySQL:
        echo !EXISTING_MYSQL_VOLUME!

        powershell -NoProfile -ExecutionPolicy Bypass -Command ^
          "$p='.env.docker'; $v='!EXISTING_MYSQL_VOLUME!'; $c=Get-Content $p -Raw; if($c -match '(?m)^MYSQL_VOLUME_NAME='){ $c=[regex]::Replace($c,'(?m)^MYSQL_VOLUME_NAME=.*$','MYSQL_VOLUME_NAME='+$v) } else { $c='MYSQL_VOLUME_NAME='+$v+[Environment]::NewLine+$c }; Set-Content -Path $p -Value $c -Encoding ASCII"
    )

    echo [INFO] Menghapus container MySQL lama tanpa menghapus volume...
    docker stop jedi-reimbursement-mysql >nul 2>nul
    docker rm jedi-reimbursement-mysql >nul 2>nul
)

for %%C in (jedi-reimbursement-backend jedi-reimbursement-frontend) do (
    docker ps -a --format "{{.Names}}" | findstr /x "%%C" >nul
    if not errorlevel 1 (
        docker stop %%C >nul 2>nul
        docker rm %%C >nul 2>nul
    )
)

echo Membangun dan menjalankan seluruh service...
docker compose --env-file .env.docker up -d --build
if errorlevel 1 (
    echo [ERROR] Full Docker gagal dijalankan.
    pause
    exit /b 1
)

echo.
echo Menunggu health check...
set /a ATTEMPT=0

:WAIT_FRONTEND
set /a ATTEMPT+=1
docker inspect --format="{{.State.Health.Status}}" jedi-reimbursement-frontend 2>nul | findstr /x "healthy" >nul
if not errorlevel 1 goto READY

if %ATTEMPT% GEQ 60 (
    echo [ERROR] Aplikasi belum sehat setelah sekitar 3 menit.
    echo Jalankan docker-logs.bat untuk melihat log.
    pause
    exit /b 1
)

timeout /t 3 /nobreak >nul
goto WAIT_FRONTEND

:READY
echo.
echo ============================================
echo Jedi Reimbursement 1.0.0 siap.
echo Buka http://localhost:!FRONTEND_PORT!
echo ============================================
pause
