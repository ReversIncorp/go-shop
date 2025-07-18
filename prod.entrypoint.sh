#!/bin/sh
set -e

echo "Running database migrations..."
goose -dir pkg/database/migrations postgres "host=db user=marketplace password=12345 dbname=go-shop sslmode=disable" up

echo "Starting application..."
./main