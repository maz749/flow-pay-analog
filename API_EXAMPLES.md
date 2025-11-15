# API Examples

Примеры использования API FlowPay

## Содержание
- [Аутентификация](#аутентификация)
- [Подписки](#подписки)
- [Telegram настройки](#telegram-настройки)

---

## Аутентификация

### Регистрация нового пользователя

**Request:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "Иван Иванов",
    "email": "ivan@example.com",
    "password": "mypassword123"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "ivan@example.com",
      "username": "Иван Иванов",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

### Вход в систему

**Request:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ivan@example.com",
    "password": "mypassword123"
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "ivan@example.com",
      "username": "Иван Иванов",
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  }
}
```

### Получение профиля

**Request:**
```bash
curl -X GET http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "ivan@example.com",
    "username": "Иван Иванов",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

---

## Подписки

### Создание подписки

**Request:**
```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Netflix",
    "description": "Стриминговый сервис",
    "amount": 999,
    "currency": "RUB",
    "billing_period": "monthly",
    "start_date": "2024-01-01",
    "category": "streaming",
    "notify_days_before": [1, 7]
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix",
    "description": "Стриминговый сервис",
    "amount": 999,
    "currency": "RUB",
    "billing_period": "monthly",
    "start_date": "2024-01-01T00:00:00Z",
    "next_billing_date": "2024-02-01T00:00:00Z",
    "is_active": true,
    "category": "streaming",
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z",
    "notifications": [
      {
        "id": 1,
        "subscription_id": 1,
        "notify_days_before": 1,
        "is_enabled": true,
        "created_at": "2024-01-15T10:00:00Z"
      },
      {
        "id": 2,
        "subscription_id": 1,
        "notify_days_before": 7,
        "is_enabled": true,
        "created_at": "2024-01-15T10:00:00Z"
      }
    ]
  }
}
```

### Создание подписки с кастомным периодом

**Request:**
```bash
curl -X POST http://localhost:8080/api/subscriptions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Хостинг",
    "amount": 500,
    "currency": "RUB",
    "billing_period": "custom",
    "custom_period_days": 90,
    "start_date": "2024-01-01",
    "category": "software",
    "notify_days_before": [3, 7]
  }'
```

### Получение всех подписок

**Request:**
```bash
curl -X GET http://localhost:8080/api/subscriptions \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Netflix",
      "description": "Стриминговый сервис",
      "amount": 999,
      "currency": "RUB",
      "billing_period": "monthly",
      "next_billing_date": "2024-02-01T00:00:00Z",
      "is_active": true,
      "category": "streaming",
      "notifications": [...]
    },
    {
      "id": 2,
      "user_id": 1,
      "name": "Spotify",
      "amount": 299,
      "currency": "RUB",
      "billing_period": "monthly",
      "next_billing_date": "2024-02-05T00:00:00Z",
      "is_active": true,
      "category": "music",
      "notifications": [...]
    }
  ]
}
```

### Получение подписки по ID

**Request:**
```bash
curl -X GET http://localhost:8080/api/subscriptions/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix",
    "description": "Стриминговый сервис",
    "amount": 999,
    "currency": "RUB",
    "billing_period": "monthly",
    "next_billing_date": "2024-02-01T00:00:00Z",
    "is_active": true,
    "category": "streaming",
    "notifications": [...]
  }
}
```

### Обновление подписки

**Request:**
```bash
curl -X PUT http://localhost:8080/api/subscriptions/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Netflix Premium",
    "amount": 1499,
    "is_active": true
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Netflix Premium",
    "amount": 1499,
    "currency": "RUB",
    "billing_period": "monthly",
    "next_billing_date": "2024-02-01T00:00:00Z",
    "is_active": true,
    "category": "streaming",
    "updated_at": "2024-01-15T11:00:00Z",
    "notifications": [...]
  }
}
```

### Деактивация подписки

**Request:**
```bash
curl -X PUT http://localhost:8080/api/subscriptions/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "is_active": false
  }'
```

### Удаление подписки

**Request:**
```bash
curl -X DELETE http://localhost:8080/api/subscriptions/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "message": "Subscription deleted successfully"
}
```

### Получение статистики

**Request:**
```bash
curl -X GET http://localhost:8080/api/subscriptions/stats \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_subscriptions": 5,
    "active_subscriptions": 4,
    "monthly_total": 3497,
    "yearly_total": 41964,
    "next_payment": {
      "id": 3,
      "name": "Spotify",
      "amount": 299,
      "currency": "RUB",
      "next_billing_date": "2024-01-20T00:00:00Z"
    }
  }
}
```

---

## Telegram настройки

### Получение настроек Telegram

**Request:**
```bash
curl -X GET http://localhost:8080/api/telegram/settings \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "chat_id": 123456789,
    "is_enabled": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  }
}
```

### Обновление настроек Telegram

**Request:**
```bash
curl -X PUT http://localhost:8080/api/telegram/settings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "chat_id": 123456789,
    "is_enabled": true
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 1,
    "chat_id": 123456789,
    "is_enabled": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Отключение уведомлений

**Request:**
```bash
curl -X PUT http://localhost:8080/api/telegram/settings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "is_enabled": false
  }'
```

---

## Примеры ошибок

### Ошибка валидации

**Response:**
```json
{
  "success": false,
  "error": "Email and password are required"
}
```

### Ошибка аутентификации

**Response:**
```json
{
  "success": false,
  "error": "Invalid or expired token"
}
```

### Подписка не найдена

**Response:**
```json
{
  "success": false,
  "error": "Subscription not found"
}
```

### Пользователь уже существует

**Response:**
```json
{
  "success": false,
  "error": "User with this email already exists"
}
```

---

## Периоды оплаты

Доступные значения для `billing_period`:

- `monthly` - Ежемесячно
- `yearly` - Ежегодно
- `weekly` - Еженедельно
- `custom` - Кастомный период (требует `custom_period_days`)

## Категории подписок

Доступные значения для `category`:

- `streaming` - Стриминг
- `music` - Музыка
- `software` - Софт
- `gaming` - Игры
- `education` - Образование
- `cloud` - Облако
- `other` - Другое

## Валюты

Доступные значения для `currency`:

- `RUB` - Российский рубль
- `USD` - Доллар США
- `EUR` - Евро
