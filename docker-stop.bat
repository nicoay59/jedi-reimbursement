@echo off
setlocal
cd /d "%~dp0"

if exist ".env.docker" (
    docker compose --env-file .env.docker down
) else (
    docker compose down
)

echo Service Docker telah dihentikan.
pause
