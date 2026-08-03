#!/bin/sh
set -eu

echo "[1/4] Menunggu database..."
/app/wait-db

echo "[2/4] Menjalankan migration..."
/app/migrate

echo "[3/4] Menjalankan seed..."
/app/seed

echo "[4/4] Menjalankan API..."
exec /app/api
