# 📱 Настройка Telegram уведомлений

## Архитектура уведомлений

**Правильный подход (для SaaS):**
- ✅ Один централизованный бот для всего приложения
- ✅ Каждый пользователь подключает свой Telegram аккаунт
- ✅ Бот отправляет персональные уведомления каждому

**Неправильный подход:**
- ❌ Каждый пользователь создает своего бота
- ❌ Пользователи вводят токены (слишком сложно)

---

## 🤖 Шаг 1: Создайте бота для приложения (один раз)

### Для владельца приложения:

1. **Откройте Telegram, найдите @BotFather**

2. **Создайте бота:**
   ```
   /newbot
   ```

3. **Введите имя бота:**
   ```
   FlowPay Notifications Bot
   ```

4. **Введите username:**
   ```
   flowpay_notifications_bot
   ```
   (Должен заканчиваться на `_bot` и быть уникальным)

5. **Сохраните токен:**
   ```
   1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
   ```

6. **Настройте бота (опционально):**
   ```
   /setdescription - Бот для уведомлений о подписках
   /setabouttext - FlowPay уведомляет о предстоящих списаниях
   /setuserpic - Загрузите логотип приложения
   ```

---

## 🔧 Шаг 2: Настройте бота на сервере

### Вариант 1: Railway/Render

В переменных окружения добавьте:
```
TELEGRAM_BOT_TOKEN=ваш-токен-бота
APP_URL=https://ваш-домен.com
```

### Вариант 2: VPS с Docker

Обновите `.env`:
```bash
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrsTUVwxyz
APP_URL=https://flowpay.example.com
```

Перезапустите:
```bash
docker-compose down
docker-compose up -d
```

---

## 👤 Шаг 3: Как пользователи подключают уведомления

### Для пользователей приложения:

1. **Откройте настройки в FlowPay**
   - Нажмите кнопку ⚙️ "Настройки"
   - Найдите раздел "Telegram уведомления"

2. **Получите код подключения**
   - Нажмите "Подключить Telegram"
   - Скопируйте уникальный код (например: `CONNECT_abc123xyz`)

3. **Откройте бота в Telegram**
   - Перейдите по ссылке: `t.me/flowpay_notifications_bot`
   - Нажмите "Start" или `/start`

4. **Отправьте код боту**
   ```
   /connect CONNECT_abc123xyz
   ```

5. **Готово!**
   - Бот подтвердит подключение
   - Теперь вы будете получать уведомления

---

## 🔔 Настройка уведомлений

После подключения можно настроить:

1. **За сколько дней уведомлять:**
   - В настройках выберите: 1, 3, 7 дней

2. **Тестовое уведомление:**
   - Нажмите "Отправить тест"
   - Бот отправит пробное сообщение

3. **Отключить уведомления:**
   - Можно временно отключить без отвязки аккаунта

---

## 💻 Технические детали (для разработчиков)

### Архитектура

```
Приложение (FlowPay)
    ↓
Telegram Bot API
    ↓
Пользователь в Telegram
```

### База данных

Таблица `telegram_settings`:
```sql
CREATE TABLE telegram_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    chat_id BIGINT UNIQUE,           -- Telegram Chat ID
    is_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP
);
```

### Процесс подключения

1. Пользователь нажимает "Подключить Telegram"
2. Бэкенд генерирует уникальный токен и сохраняет в `pending_connections`
3. Пользователь отправляет токен боту: `/connect TOKEN`
4. Бот проверяет токен, получает `user_id` и `chat_id`
5. Бот сохраняет `chat_id` в `telegram_settings`
6. Подключение завершено!

### API эндпоинты

**Получить статус подключения:**
```
GET /api/telegram/status
Authorization: Bearer {jwt_token}

Response:
{
  "connected": true,
  "chat_id": 123456789,
  "is_enabled": true
}
```

**Создать токен для подключения:**
```
POST /api/telegram/connect-token
Authorization: Bearer {jwt_token}

Response:
{
  "token": "CONNECT_abc123xyz",
  "expires_at": "2025-01-01T12:00:00Z"
}
```

**Отключить уведомления:**
```
POST /api/telegram/disconnect
Authorization: Bearer {jwt_token}
```

**Отправить тестовое уведомление:**
```
POST /api/telegram/test
Authorization: Bearer {jwt_token}
```

---

## 🤖 Код бота (Go)

### Инициализация бота

```go
// internal/telegram/bot.go
package telegram

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
)

type Bot struct {
    api *tgbotapi.BotAPI
    db  *sql.DB
}

func NewBot(token string, db *sql.DB) (*Bot, error) {
    bot, err := tgbotapi.NewBotAPI(token)
    if err != nil {
        return nil, err
    }

    bot.Debug = false
    log.Printf("Telegram bot authorized: @%s", bot.Self.UserName)

    return &Bot{
        api: bot,
        db:  db,
    }, nil
}

func (b *Bot) Start() {
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates := b.api.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        b.handleMessage(update.Message)
    }
}
```

### Обработка команд

```go
func (b *Bot) handleMessage(message *tgbotapi.Message) {
    if message.IsCommand() {
        switch message.Command() {
        case "start":
            b.handleStart(message)
        case "connect":
            b.handleConnect(message)
        case "status":
            b.handleStatus(message)
        case "help":
            b.handleHelp(message)
        }
    }
}

func (b *Bot) handleConnect(message *tgbotapi.Message) {
    args := message.CommandArguments()
    if args == "" {
        b.sendMessage(message.Chat.ID, "Используйте: /connect ВАШТОКЕН")
        return
    }

    // Проверяем токен в БД
    var userID int
    err := b.db.QueryRow(
        "SELECT user_id FROM pending_connections WHERE token = $1 AND expires_at > NOW()",
        args,
    ).Scan(&userID)

    if err != nil {
        b.sendMessage(message.Chat.ID, "❌ Неверный или истекший токен")
        return
    }

    // Сохраняем chat_id
    _, err = b.db.Exec(
        "INSERT INTO telegram_settings (user_id, chat_id, is_enabled) VALUES ($1, $2, true) ON CONFLICT (user_id) DO UPDATE SET chat_id = $2, is_enabled = true",
        userID,
        message.Chat.ID,
    )

    if err != nil {
        b.sendMessage(message.Chat.ID, "❌ Ошибка подключения")
        return
    }

    // Удаляем использованный токен
    b.db.Exec("DELETE FROM pending_connections WHERE token = $1", args)

    b.sendMessage(message.Chat.ID, "✅ Telegram успешно подключен! Теперь вы будете получать уведомления о подписках.")
}
```

### Отправка уведомлений

```go
func (b *Bot) SendNotification(chatID int64, subscription Subscription) error {
    text := fmt.Sprintf(
        "🔔 *Напоминание о подписке*\n\n"+
        "💳 %s\n"+
        "💰 %s %s\n"+
        "📅 Списание: %s\n\n"+
        "Не забудьте проверить баланс!",
        subscription.Name,
        subscription.Amount,
        subscription.Currency,
        subscription.NextBillingDate.Format("02.01.2006"),
    )

    msg := tgbotapi.NewMessage(chatID, text)
    msg.ParseMode = "Markdown"

    _, err := b.api.Send(msg)
    return err
}
```

---

## 📅 Cron-задача для уведомлений

```go
// internal/scheduler/notifications.go
package scheduler

import (
    "time"
    "log"
)

func (s *Scheduler) CheckUpcomingPayments() {
    // Запускаем каждый день в 10:00
    ticker := time.NewTicker(24 * time.Hour)
    defer ticker.Stop()

    for range ticker.C {
        s.sendNotifications()
    }
}

func (s *Scheduler) sendNotifications() {
    // Найти все подписки с платежами в ближайшие дни
    rows, err := s.db.Query(`
        SELECT
            s.id, s.name, s.amount, s.currency, s.next_billing_date,
            ts.chat_id
        FROM subscriptions s
        JOIN telegram_settings ts ON s.user_id = ts.user_id
        WHERE ts.is_enabled = true
        AND s.is_active = true
        AND s.next_billing_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'
    `)

    if err != nil {
        log.Printf("Error fetching upcoming payments: %v", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var sub Subscription
        var chatID int64

        err := rows.Scan(&sub.ID, &sub.Name, &sub.Amount, &sub.Currency, &sub.NextBillingDate, &chatID)
        if err != nil {
            continue
        }

        // Отправляем уведомление
        err = s.bot.SendNotification(chatID, sub)
        if err != nil {
            log.Printf("Error sending notification: %v", err)
        }
    }
}
```

---

## 🔐 Безопасность

### Best Practices

1. **Токен бота должен быть секретным**
   - Никогда не коммитьте в Git
   - Храните в переменных окружения
   - Используйте `.env` файл локально

2. **Токены подключения должны истекать**
   ```sql
   expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '1 hour'
   ```

3. **Один пользователь = один chat_id**
   - Используйте UNIQUE constraint на `user_id`
   - При переподключении обновляйте `chat_id`

4. **Валидация входящих сообщений**
   - Проверяйте формат токенов
   - Защищайтесь от SQL-инъекций (используйте параметризованные запросы)

---

## 🧪 Тестирование

### Локальное тестирование

1. Создайте тестового бота через @BotFather
2. Добавьте токен в `.env`
3. Запустите приложение
4. Откройте бота в Telegram
5. Получите токен подключения в UI
6. Отправьте боту: `/connect TOKEN`
7. Проверьте, что бот ответил успешно
8. Отправьте тестовое уведомление из UI

### Тестовая команда для бота

```
/start - Начать работу с ботом
/connect TOKEN - Подключить аккаунт FlowPay
/status - Проверить статус подключения
/disconnect - Отключить уведомления
/help - Справка
```

---

## 📊 Мониторинг

### Полезные SQL запросы

**Сколько пользователей подключили Telegram:**
```sql
SELECT COUNT(*) FROM telegram_settings WHERE is_enabled = true;
```

**Список активных подключений:**
```sql
SELECT u.email, ts.chat_id, ts.created_at
FROM telegram_settings ts
JOIN users u ON ts.user_id = u.id
WHERE ts.is_enabled = true;
```

**Статистика отправленных уведомлений:**
```sql
-- Если добавить таблицу notification_logs
SELECT DATE(sent_at), COUNT(*)
FROM notification_logs
WHERE status = 'sent'
GROUP BY DATE(sent_at)
ORDER BY DATE(sent_at) DESC;
```

---

## 🚀 Что нужно для продакшена

### Чеклист

- [ ] Бот создан через @BotFather
- [ ] Токен бота добавлен в переменные окружения
- [ ] UI для подключения Telegram реализован
- [ ] API эндпоинты для подключения работают
- [ ] Cron-задача для уведомлений настроена
- [ ] Тестовое уведомление работает
- [ ] Логи работы бота настроены
- [ ] Обработка ошибок добавлена

---

## 💡 Альтернативные решения

### Email уведомления

Если Telegram сложен, можно использовать email:

**Плюсы:**
- Проще реализовать
- Не требует внешних сервисов
- Email есть у всех

**Минусы:**
- Могут попадать в спам
- Менее удобно для пользователей
- Нужен SMTP сервер

### Push-уведомления в браузере

**Плюсы:**
- Встроено в браузер
- Не требует регистрации

**Минусы:**
- Работает только если сайт открыт
- Требует HTTPS
- Сложнее реализация

---

**Рекомендация:** Используйте централизованного Telegram бота как в инструкции выше. Это стандартный подход для SaaS-приложений.
