# FlowPay - Сервис управления подписками

FlowPay - это веб-приложение для управления подписками с уведомлениями в Telegram. Приложение позволяет отслеживать все ваши подписки, получать напоминания о предстоящих списаниях и контролировать расходы.

## Возможности

- 📊 **Управление подписками** - добавляйте, редактируйте и удаляйте подписки
- 💰 **Статистика расходов** - отслеживайте месячные и годовые траты
- ⏰ **Напоминания** - настройте уведомления за 1, 3, 7 или 14 дней до списания
- 📱 **Telegram интеграция** - получайте уведомления прямо в Telegram
- ❌ **Отмена подписок** - инструкции по отмене для 20+ популярных сервисов (Spotify, Netflix, Яндекс Плюс и др.)
- 🎨 **Адаптивный дизайн** - работает на всех устройствах (десктоп, планшет, мобильный)
- 🔐 **Безопасность** - JWT аутентификация и хэширование паролей

## Технологии

**Backend:**
- Go 1.21
- PostgreSQL 15
- Gorilla Mux (роутинг)
- JWT (аутентификация)
- Telegram Bot API
- Cron scheduler

**Frontend:**
- HTML5, CSS3, JavaScript (Vanilla)
- Адаптивный дизайн
- SPA архитектура

## Требования

- Docker и Docker Compose (для запуска через Docker)
- Go 1.21+ (для локальной разработки)
- PostgreSQL 15+ (для локальной разработки)
- Telegram Bot Token (опционально, для уведомлений)

## Быстрый старт с Docker

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/maz749/flow-pay-analog.git
cd flow-pay-analog
```

### 2. Настройте переменные окружения

Создайте файл `.env` в корне проекта:

```bash
cp .env.example .env
```

Отредактируйте `.env` файл и укажите свои значения:

```env
# База данных (можно оставить по умолчанию для Docker)
DB_HOST=postgres
DB_PORT=5432
DB_USER=flowpay
DB_PASSWORD=flowpay_password
DB_NAME=flowpay_db
DB_SSLMODE=disable

# Сервер
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# JWT Secret (обязательно измените в продакшене!)
JWT_SECRET=ваш-секретный-ключ-min-32-символа

# Telegram Bot (опционально)
TELEGRAM_BOT_TOKEN=ваш-токен-от-BotFather

# Окружение
APP_ENV=production
```

### 3. Создайте Telegram бота (опционально)

Если хотите использовать уведомления в Telegram:

1. Найдите [@BotFather](https://t.me/BotFather) в Telegram
2. Отправьте команду `/newbot`
3. Следуйте инструкциям и получите токен
4. Вставьте токен в `.env` файл в переменную `TELEGRAM_BOT_TOKEN`

### 4. Запустите приложение

```bash
docker-compose up -d
```

Приложение будет доступно по адресу: **http://localhost:8080**

### 5. Проверьте логи

```bash
docker-compose logs -f
```

### 6. Остановите приложение

```bash
docker-compose down
```

## Локальная разработка без Docker

### 1. Установите зависимости

```bash
go mod download
```

### 2. Настройте PostgreSQL

Создайте базу данных:

```bash
createdb flowpay_db
```

Примените миграции:

```bash
psql -d flowpay_db -f migrations/001_init_schema.sql
```

Или используйте Makefile:

```bash
make migrate
```

### 3. Настройте `.env` файл

```bash
cp .env.example .env
# Отредактируйте .env с вашими локальными настройками
```

Для локальной разработки используйте:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=ваш-пользователь
DB_PASSWORD=ваш-пароль
DB_NAME=flowpay_db
SERVER_HOST=localhost
```

### 4. Запустите приложение

```bash
make run
# или
go run cmd/server/main.go
```

### 5. Соберите бинарный файл

```bash
make build
```

Скомпилированный файл будет в `bin/flowpay`

## Использование приложения

### Регистрация и вход

1. Откройте http://localhost:8080
2. Перейдите на вкладку "Регистрация"
3. Заполните форму (имя, email, пароль)
4. После регистрации вы автоматически войдете в систему

### Добавление подписки

1. Нажмите кнопку "+ Добавить"
2. Заполните информацию о подписке:
   - Название (например, "Netflix")
   - Сумма (например, 999)
   - Валюта (RUB, USD, EUR)
   - Период оплаты (ежемесячно, ежегодно, еженедельно, другой)
   - Дата начала
   - Категория (опционально)
   - Напоминания (выберите за сколько дней уведомлять)
3. Нажмите "Сохранить"

### Настройка Telegram уведомлений

1. Нажмите на иконку настроек (⚙️) в правом верхнем углу
2. Найдите вашего Telegram бота
3. Отправьте боту команду `/start`
4. Скопируйте полученный Chat ID
5. Вставьте Chat ID в настройках приложения
6. Включите переключатель "Включить Telegram уведомления"
7. Нажмите "Сохранить"

Теперь вы будете получать уведомления о предстоящих списаниях!

### Просмотр статистики

На главной странице отображается:
- Общее количество подписок
- Сумма ежемесячных платежей
- Сумма годовых платежей
- Информация о следующем списании

## API Endpoints

### Аутентификация

```
POST /api/auth/register - Регистрация
POST /api/auth/login    - Вход
GET  /api/auth/profile  - Профиль пользователя (требует авторизации)
```

### Подписки

```
GET    /api/subscriptions        - Получить все подписки
POST   /api/subscriptions        - Создать подписку
GET    /api/subscriptions/:id    - Получить подписку по ID
PUT    /api/subscriptions/:id    - Обновить подписку
DELETE /api/subscriptions/:id    - Удалить подписку
GET    /api/subscriptions/stats  - Получить статистику
```

### Telegram

```
GET /api/telegram/settings - Получить настройки Telegram
PUT /api/telegram/settings - Обновить настройки Telegram
```

### Пример запроса

```bash
# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'

# Добавление подписки (требует токен)
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Netflix",
    "amount": 999,
    "currency": "RUB",
    "billing_period": "monthly",
    "start_date": "2024-01-01",
    "notify_days_before": [1, 7]
  }'
```

## Тестирование

### Запуск тестов

```bash
make test
# или
go test -v ./...
```

### Тестирование вручную

1. **Проверка регистрации и входа:**
   - Зарегистрируйте нового пользователя
   - Выйдите и войдите снова
   - Проверьте сохранение сессии при перезагрузке страницы

2. **Проверка управления подписками:**
   - Добавьте несколько подписок с разными периодами
   - Отредактируйте подписку
   - Удалите подписку
   - Проверьте обновление статистики

3. **Проверка Telegram интеграции:**
   - Настройте Telegram бота
   - Добавьте подписку с напоминанием на завтра
   - Проверьте получение уведомления

4. **Проверка адаптивности:**
   - Откройте приложение на мобильном устройстве
   - Проверьте все функции на разных размерах экрана

## Структура проекта

```
flow-pay-analog/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── handlers/                # HTTP handlers
│   │   ├── auth_handler.go
│   │   ├── subscription_handler.go
│   │   └── telegram_handler.go
│   ├── middleware/              # HTTP middleware
│   │   ├── auth.go
│   │   └── cors.go
│   ├── models/                  # Модели данных
│   │   ├── user.go
│   │   ├── subscription.go
│   │   └── notification.go
│   ├── repository/              # Работа с БД
│   │   ├── user_repository.go
│   │   ├── subscription_repository.go
│   │   └── notification_repository.go
│   ├── scheduler/               # Планировщик задач
│   │   └── scheduler.go
│   └── telegram/                # Telegram бот
│       └── bot.go
├── pkg/
│   ├── config/                  # Конфигурация
│   │   └── config.go
│   ├── database/                # Подключение к БД
│   │   └── database.go
│   └── utils/                   # Утилиты
│       ├── auth.go
│       ├── response.go
│       └── subscription.go
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css
│   │   └── js/
│   │       └── app.js
│   └── templates/
│       └── index.html
├── migrations/                  # SQL миграции
│   └── 001_init_schema.sql
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

## Makefile команды

```bash
make help         # Показать справку
make build        # Собрать приложение
make run          # Запустить локально
make test         # Запустить тесты
make clean        # Очистить собранные файлы
make deps         # Установить зависимости
make docker-build # Собрать Docker образ
make docker-up    # Запустить в Docker
make docker-down  # Остановить Docker
make docker-logs  # Показать логи
make migrate      # Применить миграции
```

## Troubleshooting

### Проблема: Не запускается база данных

```bash
# Проверьте статус контейнеров
docker-compose ps

# Проверьте логи PostgreSQL
docker-compose logs postgres

# Пересоздайте контейнеры
docker-compose down -v
docker-compose up -d
```

### Проблема: Ошибка подключения к БД

Убедитесь что:
- PostgreSQL запущен
- Параметры подключения в `.env` корректны
- База данных создана
- Миграции применены

### Проблема: Telegram бот не отправляет уведомления

Проверьте:
- Токен бота корректен
- Бот запущен (в логах должно быть "Telegram Bot authorized")
- Chat ID корректен (получите его через `/start` в боте)
- Уведомления включены в настройках

### Проблема: Приложение не доступно

```bash
# Проверьте, что сервер запущен
docker-compose logs app

# Проверьте занят ли порт 8080
lsof -i :8080  # Linux/Mac
netstat -ano | findstr :8080  # Windows

# Измените порт в docker-compose.yml если нужно
```

## Лицензия

MIT License

## Автор

Проект создан как аналог сервиса flow-pay.ru

## Поддержка

Если у вас возникли вопросы или проблемы, создайте issue в репозитории.
