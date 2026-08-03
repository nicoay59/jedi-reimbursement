@echo off
setlocal
cd /d "%~dp0frontend"

if not exist ".env" copy ".env.example" ".env" >nul

call npm config set registry https://registry.npmjs.org/

if not exist "node_modules\.bin\vite.cmd" (
    if exist "node_modules" rmdir /s /q node_modules
    call npm install
    if errorlevel 1 (
        echo [ERROR] npm install gagal.
        pause
        exit /b 1
    )
)

echo Frontend berjalan di http://localhost:5173
call npm run dev
pause
