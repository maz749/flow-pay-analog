#!/bin/bash

# Скрипт для применения миграции cancellation_instructions

set -e

echo "Applying cancellation instructions migration..."

# Загрузка переменных окружения
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Параметры подключения
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-flowpay}
DB_PASSWORD=${DB_PASSWORD:-flowpay_password}
DB_NAME=${DB_NAME:-flowpay_db}

echo "Database: $DB_NAME"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"

# Применение миграции
PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/002_add_cancellation_instructions.sql

echo "Migration applied successfully!"
echo "Cancellation instructions table created and populated with popular services."
