#!/bin/sh
set -e

echo "Running database migrations..."
goose -dir pkg/database/migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

echo "Starting application..."
./main