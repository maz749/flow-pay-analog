#!/bin/bash

# FlowPay Quick Start Script
# Этот скрипт быстро запустит приложение для тестирования

set -e

echo "🚀 FlowPay Quick Start"
echo "====================="
echo ""

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен!"
    echo "Установите Docker: https://docs.docker.com/get-docker/"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен!"
    echo "Установите Docker Compose: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✅ Docker найден"
echo "✅ Docker Compose найден"
echo ""

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Создаю .env файл..."
    cp .env.example .env
    echo "✅ .env создан"
else
    echo "✅ .env уже существует"
fi

echo ""
echo "🔨 Собираю и запускаю приложение..."
echo "Это может занять несколько минут при первом запуске..."
echo ""

# Stop any existing containers
docker-compose down 2>/dev/null || true

# Build and start
docker-compose up --build -d

echo ""
echo "⏳ Ожидание запуска сервисов..."
sleep 10

# Check if containers are running
if docker ps | grep -q "flowpay-app" && docker ps | grep -q "flowpay-postgres"; then
    echo ""
    echo "✅ FlowPay успешно запущен!"
    echo ""
    echo "🌐 Откройте в браузере: http://localhost:8080"
    echo ""
    echo "📊 Полезные команды:"
    echo "  - Остановить:     docker-compose down"
    echo "  - Логи:           docker-compose logs -f"
    echo "  - Перезапустить:  docker-compose restart"
    echo "  - Полная очистка: docker-compose down -v"
    echo ""
    echo "📚 Полная инструкция по деплою: DEPLOYMENT.md"
    echo ""
else
    echo ""
    echo "❌ Что-то пошло не так!"
    echo "Проверьте логи: docker-compose logs"
    exit 1
fi
