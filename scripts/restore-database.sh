#!/bin/bash

# FlowPay Database Restore Script
# Usage: ./scripts/restore-database.sh [backup_file.sql.gz]

set -e

# Configuration
CONTAINER_NAME="${CONTAINER_NAME:-flowpay-postgres}"
DB_NAME="${DB_NAME:-flowpay_db}"
DB_USER="${DB_USER:-flowpay}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check arguments
if [ $# -eq 0 ]; then
    log_error "Usage: $0 <backup_file.sql.gz>"
    log_info "Available backups:"
    ls -lh /opt/backups/flowpay/*.sql.gz 2>/dev/null || log_warn "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    log_error "Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Check if Docker container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    log_error "Container $CONTAINER_NAME is not running!"
    exit 1
fi

log_warn "⚠️  WARNING: This will REPLACE all data in database '$DB_NAME'"
log_warn "Press Ctrl+C to cancel, or wait 10 seconds to continue..."
sleep 10

log_info "Starting database restore..."
log_info "Backup file: $BACKUP_FILE"

# Decompress if needed
if [[ "$BACKUP_FILE" == *.gz ]]; then
    log_info "Decompressing backup..."
    TEMP_FILE="/tmp/flowpay_restore_$(date +%s).sql"
    gunzip -c "$BACKUP_FILE" > "$TEMP_FILE"
    RESTORE_FILE="$TEMP_FILE"
else
    RESTORE_FILE="$BACKUP_FILE"
fi

# Stop application
log_info "Stopping application..."
docker stop flowpay-app 2>/dev/null || log_warn "App container not running"

# Drop and recreate database
log_info "Recreating database..."
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "DROP DATABASE IF EXISTS ${DB_NAME};"
docker exec "$CONTAINER_NAME" psql -U "$DB_USER" -c "CREATE DATABASE ${DB_NAME};"

# Restore database
log_info "Restoring database..."
cat "$RESTORE_FILE" | docker exec -i "$CONTAINER_NAME" pg_restore -U "$DB_USER" -d "$DB_NAME" --no-owner --no-privileges

if [ $? -eq 0 ]; then
    log_info "Database restored successfully!"
else
    log_error "Database restore failed!"
    exit 1
fi

# Cleanup temp file
if [ -n "$TEMP_FILE" ]; then
    rm -f "$TEMP_FILE"
fi

# Start application
log_info "Starting application..."
docker start flowpay-app

log_info "Restore completed successfully! 🎉"
log_info "Application will be available in a few seconds..."

exit 0
