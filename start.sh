#!/bin/sh

set -e

echo "run db migration"
soucre /app/app.env
/app/migrate -path /app/migration -database "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE" -verbose up

echo "start app"
exec "$@"