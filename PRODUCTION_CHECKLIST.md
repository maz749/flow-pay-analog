# ✅ Production Checklist - Всё для идеального хостинга

## 🔒 Безопасность

### Обязательно:
- [ ] **JWT_SECRET** - минимум 32 символа, случайный
  ```bash
  openssl rand -base64 32
  ```
- [ ] **DB_PASSWORD** - сильный пароль (16+ символов)
- [ ] **HTTPS/SSL** - обязателен для продакшена
- [ ] **.env** - никогда не коммитить в Git
- [ ] **CORS** - настроить разрешенные домены
- [ ] **Rate Limiting** - защита от DDoS
- [ ] **SQL Injection** - использовать параметризованные запросы
- [ ] **XSS Protection** - санитизация ввода

### Security Headers (добавить в Nginx):
```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
add_header Content-Security-Policy "default-src 'self' https:; script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; style-src 'self' 'unsafe-inline';" always;
```

---

## 🗄️ База данных

### Production настройки:
- [ ] **DB_SSLMODE=require** - обязательно SSL
- [ ] **Connection Pooling** - ограничить количество соединений
- [ ] **Индексы** - проверить на всех часто используемых полях
- [ ] **Автобэкапы** - настроить ежедневные бэкапы
- [ ] **Репликация** - для высокой доступности (опционально)

### PostgreSQL оптимизация:
```sql
-- Добавить индексы для производительности
CREATE INDEX CONCURRENTLY idx_subscriptions_user_next_billing
ON subscriptions(user_id, next_billing_date) WHERE is_active = true;

CREATE INDEX CONCURRENTLY idx_subscriptions_active
ON subscriptions(is_active) WHERE is_active = true;

-- Автовакуум
ALTER TABLE subscriptions SET (autovacuum_vacuum_scale_factor = 0.05);
```

---

## 📊 Мониторинг и логирование

### Логирование:
- [ ] **Structured logging** - JSON формат
- [ ] **Log levels** - DEBUG, INFO, WARN, ERROR
- [ ] **Log rotation** - автоочистка старых логов
- [ ] **Centralized logs** - собирать в одном месте

### Метрики:
- [ ] **Health check endpoint** - `/health` или `/api/health`
- [ ] **Prometheus metrics** - экспорт метрик
- [ ] **Uptime monitoring** - UptimeRobot или Pingdom
- [ ] **Error tracking** - Sentry или Rollbar

### Что логировать:
```go
// Обязательно логировать:
- Все ошибки БД
- Неудачные попытки входа
- API ошибки (500, 503)
- Медленные запросы (> 1 сек)
- Telegram bot ошибки
- Важные бизнес-события (регистрация, подписка)
```

---

## 🚀 Производительность

### Кэширование:
- [ ] **Redis** - для сессий и кэша (опционально)
- [ ] **HTTP caching** - Cache-Control headers
- [ ] **Static assets** - CDN для картинок/JS/CSS
- [ ] **Database queries** - кэшировать частые запросы

### Оптимизация:
- [ ] **Gzip compression** - сжатие ответов
- [ ] **Image optimization** - сжать логотипы
- [ ] **Minify JS/CSS** - для фронтенда
- [ ] **Lazy loading** - отложенная загрузка

---

## 🔄 CI/CD и деплой

### Автоматизация:
- [ ] **GitHub Actions** - автодеплой при push
- [ ] **Docker builds** - автосборка образов
- [ ] **Health checks** - перед переключением трафика
- [ ] **Rollback plan** - откат при ошибках

### Git workflow:
```bash
main (production)
  ↓
staging (pre-production testing)
  ↓
develop (development)
  ↓
feature branches
```

---

## 📧 Email/Notifications

### Резервные уведомления:
- [ ] **Email fallback** - если Telegram недоступен
- [ ] **SMTP настройки** - для email уведомлений
- [ ] **Email templates** - красивые HTML письма
- [ ] **Unsubscribe** - возможность отписаться

### Email провайдеры:
- SendGrid (бесплатно до 100/день)
- Mailgun (бесплатно до 5000/мес)
- AWS SES (очень дешево)

---

## 🌐 Домен и DNS

### Настройка:
- [ ] **A-запись** - домен → IP сервера
- [ ] **CNAME** - www → основной домен
- [ ] **SSL сертификат** - Let's Encrypt
- [ ] **Cloudflare** - для DDoS защиты (опционально)

### DNS записи:
```
Type  Name  Value           TTL
A     @     your-server-ip  3600
A     www   your-server-ip  3600
```

---

## 🔧 Переменные окружения

### Все необходимые переменные:

```bash
# Database
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=flowpay
DB_PASSWORD=strong-password-here
DB_NAME=flowpay_db
DB_SSLMODE=require
MAX_OPEN_CONNS=25
MAX_IDLE_CONNS=5

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
ALLOWED_ORIGINS=https://flowpay.example.com

# Security
JWT_SECRET=your-32-char-secret-here
JWT_EXPIRY_HOURS=168
BCRYPT_COST=12
RATE_LIMIT=100

# Telegram
TELEGRAM_BOT_TOKEN=your-bot-token
TELEGRAM_TIMEOUT=30

# Application
APP_ENV=production
APP_URL=https://flowpay.example.com
APP_NAME=FlowPay

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Email (optional)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASSWORD=your-sendgrid-api-key
SMTP_FROM=notifications@flowpay.example.com

# Monitoring (optional)
SENTRY_DSN=your-sentry-dsn
```

---

## 🧪 Тестирование перед запуском

### Pre-production тесты:
```bash
# 1. Проверка безопасности
npm audit
go mod verify

# 2. Проверка SSL
curl -I https://flowpay.example.com

# 3. Проверка health check
curl https://flowpay.example.com/health

# 4. Load testing
ab -n 1000 -c 10 https://flowpay.example.com/

# 5. Database backup
pg_dump -U flowpay flowpay_db > backup_test.sql
```

### Функциональные тесты:
- [ ] Регистрация работает
- [ ] Вход работает
- [ ] Создание подписок работает
- [ ] Telegram подключение работает
- [ ] Уведомления отправляются
- [ ] Команды (Teams) работают
- [ ] Экспорт данных работает

---

## 🐳 Docker Production

### Оптимизированный Dockerfile:

```dockerfile
# Multi-stage build для меньшего размера
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код
COPY . .

# Собираем с оптимизацией
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o main ./cmd/server

# Production образ
FROM alpine:latest

# Добавляем CA сертификаты и timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web

# Непривилегированный пользователь
RUN adduser -D -u 1000 appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./main"]
```

### docker-compose.yml production:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    restart: always
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backups:/backups
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    # Ограничение ресурсов
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M

  app:
    build: .
    restart: always
    depends_on:
      postgres:
        condition: service_healthy
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: ${DB_USER}
      DB_PASSWORD: ${DB_PASSWORD}
      DB_NAME: ${DB_NAME}
      DB_SSLMODE: disable
      JWT_SECRET: ${JWT_SECRET}
      TELEGRAM_BOT_TOKEN: ${TELEGRAM_BOT_TOKEN}
      APP_ENV: production
    ports:
      - "8080:8080"
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M

volumes:
  postgres_data:
```

---

## 📝 Nginx конфигурация (для VPS)

### /etc/nginx/sites-available/flowpay:

```nginx
# Rate limiting
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=login_limit:10m rate=5r/m;

# Upstream
upstream flowpay_backend {
    server localhost:8080 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name flowpay.example.com www.flowpay.example.com;

    return 301 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name flowpay.example.com;

    # SSL certificates (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/flowpay.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/flowpay.example.com/privkey.pem;
    ssl_trusted_certificate /etc/letsencrypt/live/flowpay.example.com/chain.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_stapling on;
    ssl_stapling_verify on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Logging
    access_log /var/log/nginx/flowpay_access.log;
    error_log /var/log/nginx/flowpay_error.log;

    # Max body size
    client_max_body_size 10M;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss;

    # Static files
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://flowpay_backend;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # API rate limiting
    location /api/ {
        limit_req zone=api_limit burst=20 nodelay;

        proxy_pass http://flowpay_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Login rate limiting
    location /api/auth/login {
        limit_req zone=login_limit burst=3 nodelay;

        proxy_pass http://flowpay_backend;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Health check
    location /health {
        proxy_pass http://flowpay_backend;
        access_log off;
    }

    # Default location
    location / {
        proxy_pass http://flowpay_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

---

## 🔄 Автоматические бэкапы

### Скрипт бэкапа:

```bash
#!/bin/bash
# /opt/scripts/backup-flowpay.sh

set -e

BACKUP_DIR="/opt/backups/flowpay"
DATE=$(date +%Y%m%d_%H%M%S)
KEEP_DAYS=7

mkdir -p $BACKUP_DIR

# Database backup
docker exec flowpay-postgres pg_dump -U flowpay -F c flowpay_db > $BACKUP_DIR/db_$DATE.dump

# Compress
gzip $BACKUP_DIR/db_$DATE.dump

# Delete old backups
find $BACKUP_DIR -name "db_*.dump.gz" -mtime +$KEEP_DAYS -delete

# Optional: Upload to S3
# aws s3 cp $BACKUP_DIR/db_$DATE.dump.gz s3://your-bucket/backups/

echo "Backup completed: db_$DATE.dump.gz"
```

### Cron job:

```bash
# Ежедневно в 2 AM
0 2 * * * /opt/scripts/backup-flowpay.sh >> /var/log/flowpay-backup.log 2>&1
```

---

## 📈 Monitoring setup

### Health check endpoint (Go):

```go
// internal/handler/health.go
package handler

import (
    "github.com/gin-gonic/gin"
    "time"
)

func (h *Handler) HealthCheck(c *gin.Context) {
    // Check database
    var result int
    err := h.db.QueryRow("SELECT 1").Scan(&result)
    if err != nil {
        c.JSON(503, gin.H{
            "status": "unhealthy",
            "database": "down",
            "timestamp": time.Now(),
        })
        return
    }

    // Check Telegram bot (optional)
    botStatus := "unknown"
    if h.bot != nil {
        botStatus = "healthy"
    }

    c.JSON(200, gin.H{
        "status": "healthy",
        "database": "up",
        "telegram_bot": botStatus,
        "timestamp": time.Now(),
        "version": "1.0.0",
    })
}
```

### UptimeRobot setup:
1. Регистрация на uptimerobot.com (бесплатно)
2. Add New Monitor:
   - Type: HTTP(s)
   - URL: https://flowpay.example.com/health
   - Interval: 5 minutes
3. Alerts на email/Telegram

---

## 🚨 Error tracking с Sentry

### Интеграция Sentry:

```go
// cmd/server/main.go
import "github.com/getsentry/sentry-go"

func init() {
    err := sentry.Init(sentry.ClientOptions{
        Dsn: os.Getenv("SENTRY_DSN"),
        Environment: os.Getenv("APP_ENV"),
    })
    if err != nil {
        log.Fatalf("Sentry init failed: %v", err)
    }
}

// Обработка ошибок
func handleError(err error) {
    sentry.CaptureException(err)
    log.Error(err)
}
```

---

## ✅ Финальный чеклист

### Перед запуском:

**Безопасность:**
- [ ] JWT_SECRET изменен (32+ символа)
- [ ] DB_PASSWORD сильный
- [ ] HTTPS включен
- [ ] Security headers настроены
- [ ] Rate limiting включен

**База данных:**
- [ ] DB_SSLMODE=require
- [ ] Индексы созданы
- [ ] Автобэкапы настроены
- [ ] Connection pooling настроен

**Мониторинг:**
- [ ] Health check работает
- [ ] UptimeRobot настроен
- [ ] Логи пишутся
- [ ] Sentry подключен (опционально)

**Производительность:**
- [ ] Gzip включен
- [ ] Кэширование настроено
- [ ] Static assets оптимизированы

**Функциональность:**
- [ ] Все тесты пройдены
- [ ] Telegram bot работает
- [ ] Email уведомления работают (если включены)
- [ ] Команды работают

**Инфраструктура:**
- [ ] Домен настроен
- [ ] DNS записи созданы
- [ ] SSL сертификат установлен
- [ ] Nginx настроен (если VPS)
- [ ] Firewall настроен

**Документация:**
- [ ] README.md обновлен
- [ ] API документация актуальна
- [ ] Инструкции по деплою проверены

---

## 📚 Полезные команды

```bash
# Проверка всех сервисов
docker-compose ps

# Логи в реальном времени
docker-compose logs -f

# Перезапуск после обновления
docker-compose down && docker-compose up -d --build

# Проверка использования ресурсов
docker stats

# Бэкап базы
docker exec flowpay-postgres pg_dump -U flowpay flowpay_db > backup.sql

# Восстановление
docker exec -i flowpay-postgres psql -U flowpay flowpay_db < backup.sql

# Проверка SSL
curl -vI https://flowpay.example.com

# Load test
ab -n 1000 -c 10 https://flowpay.example.com/api/health
```

---

**После выполнения всех пунктов, ваше приложение готово к продакшену! 🚀**
