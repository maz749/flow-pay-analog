#!/bin/sh
set -e

echo "Starting application..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "ERROR: DATABASE_URL is not set"
    exit 1
fi

echo "Database URL is configured"

# Wait for database to be ready
echo "Waiting for database to be ready..."
for i in 1 2 3 4 5; do
    if psql "$DATABASE_URL" -c "SELECT 1" > /dev/null 2>&1; then
        echo "Database is ready!"
        break
    fi
    echo "Waiting for database... ($i/5)"
    sleep 2
done

# Run migrations
echo "Running database migrations..."
if [ -d "/root/migrations" ]; then
    for migration in /root/migrations/*.sql; do
        if [ -f "$migration" ]; then
            echo "Running migration: $(basename $migration)"
            psql "$DATABASE_URL" -f "$migration"
        fi
    done
    echo "Migrations completed successfully"
else
    echo "No migrations directory found, skipping migrations"
fi

# Start the application
echo "Starting main application..."
exec ./main
