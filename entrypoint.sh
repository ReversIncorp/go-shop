#!/bin/sh
set -e

go mod download
go install github.com/pressly/goose/v3/cmd/goose@latest

echo "Running database migrations..."
goose -dir pkg/database/migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable" up

echo "Starting application..."
go run app/main.go