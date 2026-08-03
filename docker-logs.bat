@echo off
setlocal
cd /d "%~dp0"

if exist ".env.docker" (
    docker compose --env-file .env.docker logs -f --tail=200
) else (
    docker compose logs -f --tail=200
)
