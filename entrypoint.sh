#!/bin/sh
set -e

# Wait for PostgreSQL
until pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}"; do
  echo "Waiting for PostgreSQL..."
  sleep 1
done

echo "Running database migrations..."
goose -dir pkg/database/migrations postgres "host=${DB_HOST} port=${DB_PORT} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

echo "Starting application..."
./main