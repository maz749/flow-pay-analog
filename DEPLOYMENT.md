# 🚀 Инструкция по хостингу и деплою FlowPay

## 📋 Содержание
1. [Локальное тестирование](#локальное-тестирование)
2. [Подготовка к продакшену](#подготовка-к-продакшену)
3. [Варианты хостинга](#варианты-хостинга)
4. [Деплой на VPS/VDS](#деплой-на-vpsvds)
5. [Деплой на Railway](#деплой-на-railway)
6. [Деплой на Render](#деплой-на-render)
7. [Деплой на DigitalOcean App Platform](#деплой-на-digitalocean)
8. [Настройка домена](#настройка-домена)
9. [Мониторинг и обслуживание](#мониторинг)

---

## 🧪 Локальное тестирование

### Вариант 1: С помощью Docker (рекомендуется)

**Требования:**
- Docker и Docker Compose установлены

**Шаги:**

```bash
# 1. Клонируйте репозиторий (если еще не сделано)
cd /path/to/subscriptionmanager

# 2. Создайте .env файл
cp .env.example .env

# 3. Отредактируйте .env (опционально)
nano .env

# 4. Запустите приложение
docker-compose up --build

# 5. Откройте в браузере
# http://localhost:8080
```

**Проверка работы:**
- ✅ Приложение открывается на http://localhost:8080
- ✅ Можно зарегистрироваться
- ✅ Можно войти
- ✅ Можно добавлять подписки
- ✅ Графики отображаются
- ✅ Команды работают (если настроено)

**Остановка:**
```bash
docker-compose down
```

**Полная очистка (включая БД):**
```bash
docker-compose down -v
```

### Вариант 2: Без Docker

**Требования:**
- Go 1.21+
- PostgreSQL 15+

**Шаги:**

```bash
# 1. Установите PostgreSQL
# Ubuntu/Debian:
sudo apt update
sudo apt install postgresql postgresql-contrib

# MacOS:
brew install postgresql@15

# 2. Создайте базу данных
sudo -u postgres psql
CREATE DATABASE flowpay_db;
CREATE USER flowpay WITH PASSWORD 'flowpay_password';
GRANT ALL PRIVILEGES ON DATABASE flowpay_db TO flowpay;
\q

# 3. Примените миграции
psql -U flowpay -d flowpay_db -f migrations/001_init_schema.sql
psql -U flowpay -d flowpay_db -f migrations/002_add_cancellation_instructions.sql
psql -U flowpay -d flowpay_db -f migrations/003_add_teams_and_roles.sql

# 4. Создайте .env файл
cp .env.example .env

# 5. Соберите и запустите
go mod download
go build -o bin/server ./cmd/server
./bin/server

# Или просто:
go run cmd/server/main.go
```

---

## 🔒 Подготовка к продакшену

### 1. Создайте production .env файл

```bash
# .env.production
DB_HOST=your-db-host
DB_PORT=5432
DB_USER=flowpay
DB_PASSWORD=STRONG_PASSWORD_HERE
DB_NAME=flowpay_db
DB_SSLMODE=require

SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# ВАЖНО: Сгенерируйте сильный секретный ключ!
JWT_SECRET=your-very-strong-secret-key-min-32-chars-12345678

# Telegram (опционально)
TELEGRAM_BOT_TOKEN=your-bot-token

APP_ENV=production
```

**Генерация безопасного JWT_SECRET:**
```bash
# Linux/Mac:
openssl rand -base64 32

# Или через Go:
go run -c 'import "crypto/rand"; import "encoding/base64"; b := make([]byte, 32); rand.Read(b); println(base64.StdEncoding.EncodeToString(b))'
```

### 2. Проверьте безопасность

**Чек-лист безопасности:**
- [ ] JWT_SECRET изменен и достаточно сложный (минимум 32 символа)
- [ ] DB_PASSWORD сильный (минимум 16 символов, буквы+цифры+символы)
- [ ] DB_SSLMODE установлен в `require` для продакшена
- [ ] .env файл добавлен в .gitignore
- [ ] Не используются дефолтные пароли

---

## 🌐 Варианты хостинга

### Сравнение платформ

| Платформа | Сложность | Цена/мес | PostgreSQL | Автодеплой |
|-----------|-----------|----------|------------|------------|
| **Railway** | ⭐ Легко | $5-20 | ✅ Включено | ✅ Да |
| **Render** | ⭐ Легко | $7-25 | ✅ Включено | ✅ Да |
| **DigitalOcean App Platform** | ⭐⭐ Средне | $5-25 | ✅ Включено | ✅ Да |
| **VPS (DigitalOcean/Linode)** | ⭐⭐⭐ Сложно | $6-20 | ❌ Настроить | ❌ Настроить |
| **Heroku** | ⭐ Легко | $7-25 | ✅ Включено | ✅ Да |

**Рекомендация:** Railway или Render для начинающих, VPS для опытных.

---

## 🚂 Деплой на Railway

**Плюсы:**
- Самый простой способ
- Бесплатный tier ($5 кредитов)
- Автоматический деплой из Git
- PostgreSQL включен

**Шаги:**

1. **Регистрация:**
   - Перейдите на https://railway.app
   - Войдите через GitHub

2. **Создайте новый проект:**
   - Нажмите "New Project"
   - Выберите "Deploy from GitHub repo"
   - Выберите ваш репозиторий `subscriptionmanager`

3. **Добавьте PostgreSQL:**
   - Нажмите "+ New"
   - Выберите "Database" → "PostgreSQL"
   - Railway автоматически создаст БД

4. **Настройте переменные окружения:**

   Перейдите в настройки вашего приложения → Variables:

   ```
   DATABASE_URL=${{Postgres.DATABASE_URL}}
   JWT_SECRET=your-strong-secret-here
   APP_ENV=production
   PORT=8080
   ```

5. **Деплой:**
   - Railway автоматически соберет и задеплоит
   - Получите публичный URL в разделе "Settings" → "Domains"

6. **Примените миграции:**

   В Railway Console выполните:
   ```bash
   psql $DATABASE_URL < /app/migrations/001_init_schema.sql
   psql $DATABASE_URL < /app/migrations/002_add_cancellation_instructions.sql
   psql $DATABASE_URL < /app/migrations/003_add_teams_and_roles.sql
   ```

**Цена:** ~$5-10/мес

---

## 🎨 Деплой на Render

**Плюсы:**
- Бесплатный tier (с ограничениями)
- Простой интерфейс
- Автоматический SSL

**Шаги:**

1. **Регистрация:**
   - https://render.com
   - Войдите через GitHub

2. **Создайте PostgreSQL базу:**
   - Dashboard → "New" → "PostgreSQL"
   - Выберите регион
   - Выберите план (Free или Starter $7/мес)
   - Сохраните Internal Database URL

3. **Создайте Web Service:**
   - Dashboard → "New" → "Web Service"
   - Подключите GitHub репозиторий
   - Настройки:
     - **Name:** flowpay
     - **Environment:** Docker
     - **Region:** выберите ближайший
     - **Branch:** main (или ваша ветка)
     - **Instance Type:** Free или Starter

4. **Переменные окружения:**

   ```
   DATABASE_URL=postgres://...  (из шага 2)
   JWT_SECRET=your-strong-secret
   APP_ENV=production
   PORT=8080
   ```

5. **Deploy:**
   - Нажмите "Create Web Service"
   - Render автоматически соберет Docker image

6. **Примените миграции:**
   - В Render Shell выполните миграции

**Цена:** $0 (Free tier) или $7-25/мес

---

## 💧 Деплой на DigitalOcean App Platform

**Плюсы:**
- Хорошая производительность
- Простая настройка
- Интеграция с DO инфраструктурой

**Шаги:**

1. **Создайте аккаунт:**
   - https://digitalocean.com
   - Получите $200 кредитов (для новых пользователей)

2. **Создайте Managed PostgreSQL:**
   - Create → Databases → PostgreSQL
   - Выберите план (Basic $15/мес)
   - Сохраните connection details

3. **Создайте App:**
   - Create → Apps → GitHub
   - Выберите репозиторий
   - DigitalOcean автоматически определит Dockerfile

4. **Настройте Environment Variables:**
   ```
   DB_HOST=your-db-host.db.ondigitalocean.com
   DB_PORT=25060
   DB_USER=doadmin
   DB_PASSWORD=your-password
   DB_NAME=defaultdb
   DB_SSLMODE=require
   JWT_SECRET=your-secret
   APP_ENV=production
   ```

5. **Deploy:**
   - Нажмите "Next" → "Create Resources"

**Цена:** ~$12-30/мес (App $5 + Database $15)

---

## 🖥️ Деплой на VPS/VDS

**Для опытных пользователей.**

**Выбор провайдера:**
- DigitalOcean Droplet ($6/мес)
- Linode ($5/мес)
- Vultr ($6/мес)
- Hetzner ($5/мес) - самый дешевый

**Требования VPS:**
- 1GB RAM (минимум)
- 25GB SSD
- Ubuntu 22.04 LTS

### Пошаговая инструкция:

**1. Создайте VPS и подключитесь:**

```bash
ssh root@your-server-ip
```

**2. Установите Docker и Docker Compose:**

```bash
# Обновите систему
apt update && apt upgrade -y

# Установите Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Установите Docker Compose
apt install docker-compose -y

# Проверьте установку
docker --version
docker-compose --version
```

**3. Установите Nginx (опционально, для SSL):**

```bash
apt install nginx certbot python3-certbot-nginx -y
```

**4. Клонируйте репозиторий:**

```bash
cd /opt
git clone https://github.com/your-username/subscriptionmanager.git
cd subscriptionmanager
```

**5. Создайте production .env:**

```bash
cp .env.example .env
nano .env
```

Заполните настоящими значениями (см. раздел "Подготовка к продакшену")

**6. Запустите приложение:**

```bash
docker-compose up -d
```

**7. Настройте Nginx (для доменного имени):**

```bash
nano /etc/nginx/sites-available/flowpay
```

Вставьте:
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# Активируйте конфиг
ln -s /etc/nginx/sites-available/flowpay /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

**8. Настройте SSL (Let's Encrypt):**

```bash
certbot --nginx -d your-domain.com
```

**9. Настройте автозапуск:**

```bash
# Docker Compose уже настроит автозапуск контейнеров
# Убедитесь, что Docker запускается при загрузке:
systemctl enable docker
```

**10. Настройте firewall:**

```bash
ufw allow 22/tcp
ufw allow 80/tcp
ufw allow 443/tcp
ufw enable
```

**Обновление приложения:**

```bash
cd /opt/subscriptionmanager
git pull
docker-compose down
docker-compose up -d --build
```

---

## 🌍 Настройка домена

### Купите домен

**Регистраторы:**
- Namecheap (~$10/год)
- Google Domains (~$12/год)
- Cloudflare (~$10/год)

### Настройте DNS

Добавьте A-запись, указывающую на ваш сервер:

```
Type: A
Name: @
Value: your-server-ip
TTL: 3600
```

Для поддомена:
```
Type: A
Name: flowpay
Value: your-server-ip
TTL: 3600
```

### Проверка распространения DNS

```bash
# Linux/Mac:
dig your-domain.com

# Или используйте онлайн:
# https://www.whatsmydns.net/
```

---

## 📊 Мониторинг и обслуживание

### Проверка логов

**Docker:**
```bash
docker-compose logs -f app
docker-compose logs -f postgres
```

**VPS без Docker:**
```bash
journalctl -u flowpay -f
```

### Бэкап базы данных

**Автоматический бэкап (cron):**

```bash
# Создайте скрипт бэкапа
nano /opt/backup.sh
```

```bash
#!/bin/bash
BACKUP_DIR="/opt/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

docker exec flowpay-postgres pg_dump -U flowpay flowpay_db > $BACKUP_DIR/backup_$DATE.sql

# Удалить старые бэкапы (старше 7 дней)
find $BACKUP_DIR -name "backup_*.sql" -mtime +7 -delete
```

```bash
chmod +x /opt/backup.sh

# Добавьте в cron (каждый день в 2 AM)
crontab -e
```

Добавьте:
```
0 2 * * * /opt/backup.sh
```

### Мониторинг производительности

**Установите monitoring stack (опционально):**

```bash
# Используйте Prometheus + Grafana
# Или простой мониторинг:
apt install htop iotop
```

### Обновление приложения

**Railway/Render:**
- Просто push в Git - автодеплой

**VPS:**
```bash
cd /opt/subscriptionmanager
git pull origin main
docker-compose down
docker-compose up -d --build
```

---

## ✅ Финальный чеклист

Перед запуском в продакшн:

- [ ] JWT_SECRET изменен на сильный
- [ ] Пароли БД изменены
- [ ] SSL сертификат настроен (HTTPS)
- [ ] Бэкапы настроены
- [ ] Firewall настроен
- [ ] Домен подключен
- [ ] Логи проверены
- [ ] Приложение протестировано:
  - [ ] Регистрация работает
  - [ ] Вход работает
  - [ ] Подписки создаются
  - [ ] Графики отображаются
  - [ ] Команды работают
- [ ] Мониторинг настроен

---

## 🆘 Troubleshooting

### Проблема: "Cannot connect to database"

**Решение:**
```bash
# Проверьте, запущен ли PostgreSQL
docker ps | grep postgres

# Проверьте логи БД
docker-compose logs postgres

# Проверьте переменные окружения
docker-compose config
```

### Проблема: "Port 8080 already in use"

**Решение:**
```bash
# Найдите процесс
lsof -i :8080

# Или измените порт в docker-compose.yml
ports:
  - "8081:8080"
```

### Проблема: "502 Bad Gateway" в Nginx

**Решение:**
```bash
# Проверьте, запущено ли приложение
docker ps

# Проверьте логи Nginx
tail -f /var/log/nginx/error.log

# Проверьте SELinux (если используется)
setsebool -P httpd_can_network_connect 1
```

---

## 📚 Дополнительные ресурсы

- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Let's Encrypt Documentation](https://letsencrypt.org/docs/)
- [Nginx Documentation](https://nginx.org/en/docs/)

---

**Удачи с деплоем! 🚀**
