# Настройка Telegram авторизации

## Обзор

FlowPay теперь поддерживает авторизацию через Telegram! Пользователи могут входить используя свой Telegram аккаунт, и автоматически получать уведомления о подписках в Telegram.

## Что нужно сделать

### 1. Создать или настроить Telegram бота

Если у вас еще нет бота, создайте его через [@BotFather](https://t.me/BotFather):

1. Откройте Telegram и найдите [@BotFather](https://t.me/BotFather)
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Сохраните токен бота (например: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`)

### 2. Настроить домен для Telegram Login Widget

Telegram Login Widget требует, чтобы ваше приложение было доступно по домену (не localhost).

Выполните команду в BotFather:
```
/setdomain
```

Выберите вашего бота и укажите домен, например:
- `example.com`
- `app.example.com`
- `flowpay.herokuapp.com`

**Важно:** Для разработки на localhost вы можете использовать:
- ngrok или локальный туннель
- Добавить домен в `/etc/hosts` (не рекомендуется для продакшена)

### 3. Обновить переменные окружения

Добавьте или обновите следующие переменные в вашем `.env` файле:

```env
# Telegram Bot Configuration
TELEGRAM_BOT_TOKEN=ваш_токен_бота_здесь

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=flowpay

# JWT Secret
JWT_SECRET=your_jwt_secret_key

# Server Configuration
SERVER_HOST=0.0.0.0
PORT=8080
APP_ENV=development
```

### 4. Обновить имя бота в JavaScript

Откройте `web/static/js/app.js` и найдите функцию `initTelegramWidget`:

```javascript
function initTelegramWidget(elementId) {
    // ...
    script.setAttribute('data-telegram-login', 'YOUR_BOT_USERNAME'); // Замените на имя вашего бота
    // ...
}
```

Замените `YOUR_BOT_USERNAME` на имя вашего бота (без @), например:
```javascript
script.setAttribute('data-telegram-login', 'FlowPayBot');
```

### 5. Применить миграцию базы данных

Примените новую миграцию для добавления полей Telegram авторизации:

```bash
# Подключитесь к PostgreSQL
psql -U postgres -d flowpay

# Выполните миграцию
\i migrations/002_add_telegram_auth.sql

# Проверьте изменения
\d users
```

Или используйте ваш инструмент миграций:
```bash
# Пример с golang-migrate
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/flowpay?sslmode=disable" up
```

### 6. Перезапустить приложение

```bash
# Остановите текущий процесс (Ctrl+C)

# Перезапустите приложение
go run cmd/server/main.go
```

## Как это работает

### Для пользователей

1. Пользователь открывает FlowPay
2. Нажимает "Вход" или "Регистрация"
3. Видит кнопку "Login with Telegram"
4. Кликает на кнопку и авторизуется через Telegram
5. Автоматически создается аккаунт (если новый пользователь)
6. Автоматически настраиваются Telegram уведомления
7. Получает уведомления о списаниях прямо в Telegram!

### Технические детали

**Регистрация через Telegram:**
- Создается пользователь с `auth_type = 'telegram'`
- Сохраняется Telegram ID, username, имя и фото
- Автоматически создается `telegram_settings` с `chat_id = telegram_user_id`
- Уведомления включены по умолчанию

**Уведомления:**
- Бот отправляет сообщения напрямую пользователю (chat_id = telegram_user_id)
- Не требуется отдельная настройка Chat ID
- Пользователь сразу получает уведомления после регистрации

**Безопасность:**
- Telegram данные валидируются через HMAC-SHA256
- Проверяется подпись с использованием bot token
- Данные действительны только 24 часа

## Совместимость

Существующие пользователи с email/password авторизацией продолжат работать как раньше.

Новые пользователи могут выбрать:
- Регистрация через Email
- Регистрация через Telegram

## Тестирование

1. Откройте приложение в браузере
2. Нажмите "Вход"
3. Убедитесь, что кнопка "Login with Telegram" отображается
4. Нажмите на кнопку
5. Авторизуйтесь через Telegram
6. Проверьте, что вы автоматически вошли в приложение
7. Добавьте подписку
8. Проверьте, что уведомление приходит в Telegram

## Устранение неполадок

### Кнопка Telegram не отображается

Проверьте:
- Telegram Widget script загружен в `<head>` секции
- Правильно указан `data-telegram-login` в JavaScript
- Домен настроен в BotFather

### Ошибка "Invalid hash"

Проверьте:
- `TELEGRAM_BOT_TOKEN` в `.env` файле совпадает с токеном вашего бота
- Данные от Telegram не старше 24 часов

### Уведомления не приходят

Проверьте:
- Бот запущен (в логах должно быть "Bot started successfully")
- `telegram_settings` созданы для пользователя
- `is_enabled = true` в `telegram_settings`
- Scheduler запущен и работает

## Дополнительная информация

- [Telegram Login Widget Documentation](https://core.telegram.org/widgets/login)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [BotFather Commands](https://core.telegram.org/bots#botfather)
