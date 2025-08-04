#!/bin/sh

set -e

export DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/?sslmode=disable"

echo "Running migrations..."
goose -dir /migrations postgres "$DB_URL" up

echo "Migrations completed"