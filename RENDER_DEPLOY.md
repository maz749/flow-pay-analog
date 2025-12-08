# Деплой на Render.com

Подробная инструкция по развертыванию Subscription Manager на Render.com.

## Шаг 1: Подготовка репозитория

Убедитесь что все изменения запушены в GitHub:

```bash
git status
git push origin claude/fix-railway-db-connection-01WpCEtADcVTCyyryK1NrH2Q
```

## Шаг 2: Создание аккаунта на Render

1. Перейдите на [render.com](https://render.com)
2. Нажмите **"Get Started"** или **"Sign Up"**
3. Войдите через **GitHub** (рекомендуется для автоматических деплоев)
4. Авторизуйте Render доступ к вашим репозиториям

## Шаг 3: Создание PostgreSQL базы данных

1. На Dashboard нажмите **"New +"** → **"PostgreSQL"**
2. Заполните форму:
   - **Name**: `subscriptionmanager-db` (или любое имя)
   - **Database**: `flowpay_db`
   - **User**: `flowpay` (автоматически)
   - **Region**: выберите ближайший регион (например, Frankfurt для Европы)
   - **PostgreSQL Version**: выберите последнюю (14 или 15)
   - **Plan**: **Free** (бесплатный tier)
3. Нажмите **"Create Database"**
4. Дождитесь создания БД (может занять 1-2 минуты)
5. **ВАЖНО**: Скопируйте **"Internal Database URL"** - она понадобится для приложения

## Шаг 4: Создание Web Service

1. На Dashboard нажмите **"New +"** → **"Web Service"**
2. Выберите **"Build and deploy from a Git repository"** → **"Next"**
3. Найдите ваш репозиторий `subscriptionmanager` → нажмите **"Connect"**
4. Заполните форму:

### Основные настройки:

- **Name**: `subscriptionmanager` (будет в URL)
- **Region**: тот же регион что и база данных
- **Branch**: `claude/fix-railway-db-connection-01WpCEtADcVTCyyryK1NrH2Q`
- **Runtime**: **Docker** (Render автоматически определит Dockerfile)
- **Instance Type**: **Free** (бесплатный tier)

### Дополнительные настройки (Advanced):

- **Dockerfile Path**: `./Dockerfile` (должно определиться автоматически)
- **Docker Context Directory**: `.` (корень репозитория)
- **Auto-Deploy**: **Yes** (автоматический деплой при push в GitHub)

5. Нажмите **"Create Web Service"** (пока НЕ нажимайте, сначала добавим переменные!)

## Шаг 5: Настройка переменных окружения

**ПЕРЕД созданием сервиса** прокрутите вниз до секции **"Environment Variables"** и добавьте:

### Обязательные переменные:

| Key | Value |
|-----|-------|
| `DATABASE_URL` | Вставьте **Internal Database URL** из шага 3 (должен начинаться с `postgresql://`) |
| `JWT_SECRET` | Любая случайная строка минимум 32 символа (например: `da2556559e7bf2a2da4becd19dc3de0a`) |
| `APP_ENV` | `production` |

### Опциональные переменные:

| Key | Value |
|-----|-------|
| `TELEGRAM_BOT_TOKEN` | Токен вашего Telegram бота (если нужен) |

**ВАЖНО**: Render автоматически устанавливает переменную `PORT` - не добавляйте её вручную!

6. После добавления переменных нажмите **"Create Web Service"**

## Шаг 6: Ожидание деплоя

1. Render начнет автоматический build и deploy
2. Откройте вкладку **"Logs"** чтобы видеть процесс
3. Build может занять 3-5 минут

### Что должно произойти:

```
==> Building...
==> Dockerfile detected
==> Building image...
==> Running build...

==> Deploying...
Starting application...
Database URL is configured
Waiting for database to be ready...
Database is ready!
Running migration: 001_init_schema.sql
Migrations completed successfully
Loading configuration...
✓ DATABASE_URL found, parsing connection string...
✓ Successfully parsed DATABASE_URL: host=xxx port=5432 dbname=xxx sslmode=require
Successfully connected to database
Server starting on 0.0.0.0:10000
```

4. Когда увидите **"Your service is live 🎉"** - приложение готово!

## Шаг 7: Получение URL приложения

1. На странице вашего сервиса вверху вы увидите URL вида:
   - `https://subscriptionmanager-xxxx.onrender.com`
2. Нажмите на URL чтобы открыть ваше приложение
3. Вы должны увидеть главную страницу приложения

## Шаг 8: Проверка работоспособности

Проверьте что API работает:

```bash
# Регистрация нового пользователя
curl -X POST https://your-app.onrender.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser"
  }'
```

Должен вернуться токен JWT.

## Настройка домена (опционально)

1. Перейдите в **"Settings"** → **"Custom Domain"**
2. Нажмите **"Add Custom Domain"**
3. Следуйте инструкциям для настройки DNS

## Автоматические деплои

Render автоматически делает редеплой при каждом push в выбранную ветку GitHub.

Чтобы обновить приложение:
1. Сделайте изменения в коде
2. Закоммитьте и запушьте в GitHub
3. Render автоматически начнет новый deploy

## Мониторинг и логи

### Просмотр логов:
1. Откройте ваш сервис в Dashboard
2. Перейдите на вкладку **"Logs"**
3. Логи в реальном времени

### Метрики:
1. Вкладка **"Metrics"** показывает:
   - CPU usage
   - Memory usage
   - Request rate
   - Response time

## Troubleshooting

### Приложение не запускается

**Проверьте логи на ошибки:**
1. Откройте **"Logs"**
2. Найдите строки с `✗` или `ERROR`
3. Проверьте что `DATABASE_URL` правильный

### База данных недоступна

**Проверьте:**
1. PostgreSQL сервис работает (статус "Available")
2. Используете **Internal Database URL** (не External!)
3. Приложение и БД в одном регионе

### Медленный первый запрос

Бесплатный tier Render засыпает после 15 минут неактивности. Первый запрос может занять 30-60 секунд (cold start).

**Решение**: Upgrade на платный план ($7/месяц) для постоянной работы.

## Сравнение Free tier Render vs Railway

| Функция | Render Free | Railway Free |
|---------|-------------|--------------|
| Часов в месяц | 750 часов | 500 часов |
| RAM | 512 MB | 512 MB |
| CPU | Shared | Shared |
| Databases | 1 PostgreSQL (90 дней) | Unlimited |
| Auto-sleep | После 15 мин | Нет |
| Build time | До 10 мин | До 30 мин |

**Важно**: PostgreSQL на Render Free tier удаляется через 90 дней. После этого нужно будет пересоздать БД или перейти на платный план ($7/месяц).

## Следующие шаги

После успешного деплоя:
1. ✅ Настройте Telegram бота (если нужно)
2. ✅ Подключите фронтенд к вашему API
3. ✅ Настройте CORS если нужно
4. ✅ Добавьте свой домен

## Полезные ссылки

- [Render Documentation](https://render.com/docs)
- [Render Community](https://community.render.com)
- [Render Status](https://status.render.com)

---

**Вопросы?** Проверьте логи в Render Dashboard или обратитесь в поддержку Render.
