#!/bin/sh
set -e

echo "Starting application..."

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "ERROR: DATABASE_URL is not set"
    exit 1
fi

echo "Database URL is configured"

# Parse DATABASE_URL to get connection parameters for psql
# Format: postgresql://user:password@host:port/database
DB_USER=$(echo $DATABASE_URL | sed -n 's/.*:\/\/\([^:]*\):.*/\1/p')
DB_PASS=$(echo $DATABASE_URL | sed -n 's/.*:\/\/[^:]*:\([^@]*\)@.*/\1/p')
DB_HOST=$(echo $DATABASE_URL | sed -n 's/.*@\([^:]*\):.*/\1/p')
DB_PORT=$(echo $DATABASE_URL | sed -n 's/.*:\([0-9]*\)\/.*/\1/p')
DB_NAME=$(echo $DATABASE_URL | sed -n 's/.*\/\([^?]*\).*/\1/p')

# Wait for database to be ready
echo "Waiting for database to be ready..."
for i in 1 2 3 4 5; do
    if PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1" > /dev/null 2>&1; then
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
            PGPASSWORD=$DB_PASS psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
        fi
    done
    echo "Migrations completed successfully"
else
    echo "No migrations directory found, skipping migrations"
fi

# Start the application
echo "Starting main application..."
exec ./main
