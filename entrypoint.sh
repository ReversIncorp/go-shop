#!/bin/sh
set -e

go mod download
go install github.com/pressly/goose/v3/cmd/goose@latest

echo "Running database migrations..."
goose -dir pkg/database/migrations postgres "host=db user=marketplace password=12345 dbname=go-shop sslmode=disable" up

echo "Starting application..."
go run app/main.go 