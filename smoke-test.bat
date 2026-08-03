@echo off
setlocal EnableExtensions
cd /d "%~dp0"

set FRONTEND_PORT=8088
if exist ".env.docker" (
    for /f "usebackq tokens=1,* delims==" %%A in (".env.docker") do (
        if /I "%%A"=="FRONTEND_PORT" set FRONTEND_PORT=%%B
    )
)

powershell -NoProfile -ExecutionPolicy Bypass ^
-File "%~dp0smoke-test.ps1" ^
-FrontendUrl "http://localhost:%FRONTEND_PORT%" ^
-ApiBaseUrl "http://localhost:%FRONTEND_PORT%/api/v1"

if errorlevel 1 (
    echo.
    echo Periksa log dengan docker-logs.bat.
    pause
    exit /b 1
)

pause
