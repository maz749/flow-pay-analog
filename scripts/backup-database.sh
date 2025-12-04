#!/bin/bash

# FlowPay Database Backup Script
# Usage: ./scripts/backup-database.sh

set -e

# Configuration
BACKUP_DIR="${BACKUP_DIR:-/opt/backups/flowpay}"
KEEP_DAYS="${KEEP_DAYS:-7}"
DATE=$(date +%Y%m%d_%H%M%S)
CONTAINER_NAME="${CONTAINER_NAME:-flowpay-postgres}"
DB_NAME="${DB_NAME:-flowpay_db}"
DB_USER="${DB_USER:-flowpay}"

# S3 Configuration (optional)
S3_BUCKET="${S3_BUCKET:-}"
S3_PATH="${S3_PATH:-backups/flowpay}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker container is running
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    log_error "Container $CONTAINER_NAME is not running!"
    exit 1
fi

# Create backup directory
mkdir -p "$BACKUP_DIR"

log_info "Starting database backup..."
log_info "Date: $DATE"
log_info "Container: $CONTAINER_NAME"
log_info "Database: $DB_NAME"

# Backup database
BACKUP_FILE="$BACKUP_DIR/flowpay_${DATE}.sql"
log_info "Backing up to: $BACKUP_FILE"

docker exec "$CONTAINER_NAME" pg_dump -U "$DB_USER" -F c "$DB_NAME" > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    log_info "Database backup completed successfully"
else
    log_error "Database backup failed!"
    exit 1
fi

# Compress backup
log_info "Compressing backup..."
gzip "$BACKUP_FILE"
BACKUP_FILE_GZ="${BACKUP_FILE}.gz"

if [ $? -eq 0 ]; then
    log_info "Compression completed: $BACKUP_FILE_GZ"

    # Get file size
    SIZE=$(du -h "$BACKUP_FILE_GZ" | cut -f1)
    log_info "Backup size: $SIZE"
else
    log_error "Compression failed!"
    exit 1
fi

# Upload to S3 (if configured)
if [ -n "$S3_BUCKET" ]; then
    log_info "Uploading to S3..."

    if command -v aws &> /dev/null; then
        aws s3 cp "$BACKUP_FILE_GZ" "s3://${S3_BUCKET}/${S3_PATH}/flowpay_${DATE}.sql.gz"

        if [ $? -eq 0 ]; then
            log_info "Uploaded to S3: s3://${S3_BUCKET}/${S3_PATH}/flowpay_${DATE}.sql.gz"
        else
            log_warn "S3 upload failed (backup still saved locally)"
        fi
    else
        log_warn "AWS CLI not found, skipping S3 upload"
    fi
fi

# Delete old backups
log_info "Cleaning up old backups (older than $KEEP_DAYS days)..."
find "$BACKUP_DIR" -name "flowpay_*.sql.gz" -mtime +$KEEP_DAYS -delete

REMAINING=$(find "$BACKUP_DIR" -name "flowpay_*.sql.gz" | wc -l)
log_info "Backups remaining: $REMAINING"

# Create latest symlink
ln -sf "$BACKUP_FILE_GZ" "$BACKUP_DIR/flowpay_latest.sql.gz"

log_info "Backup completed successfully! 🎉"
log_info "Backup location: $BACKUP_FILE_GZ"

# Optional: Send notification
# curl -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
#      -d "chat_id=${ADMIN_CHAT_ID}" \
#      -d "text=✅ FlowPay backup completed: $SIZE"

exit 0
