@echo off
setlocal
cd /d "%~dp0"

call "%~dp0start-database.bat"
if errorlevel 1 exit /b 1

start "Jedi Reimbursement Backend" cmd /k call "%~dp0run-backend.bat"
timeout /t 4 /nobreak >nul
start "Jedi Reimbursement Frontend" cmd /k call "%~dp0run-frontend.bat"

echo Backend dan frontend dijalankan pada jendela terpisah.
echo Buka http://localhost:5173
