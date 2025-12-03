# Функциональность отмены подписок

## Обзор

Новая функциональность позволяет пользователям легко находить инструкции по отмене подписок для популярных сервисов.

## Компоненты

### 1. База данных

**Таблица `cancellation_instructions`:**
- `id` - уникальный идентификатор
- `service_name` - название сервиса (например, "Spotify", "Netflix")
- `instructions` - подробные инструкции по отмене
- `note` - дополнительные примечания
- `url` - ссылка на страницу отмены подписки
- `category` - категория сервиса

**Миграция:**
```bash
# Применить миграцию
./scripts/migrate-cancellation.sh

# Или через Make
make migrate
```

### 2. Backend API

**Endpoints:**

1. **GET `/api/cancellation/instructions`** - Получить все инструкции
   - Query параметр `?search=query` для поиска

2. **GET `/api/cancellation/instructions/{serviceName}`** - Получить инструкцию для конкретного сервиса

3. **POST `/api/subscriptions/{id}/cancel`** - Отменить подписку пользователя
   - Требует аутентификацию
   - Помечает подписку как неактивную

**Пример запроса:**
```bash
# Получить все инструкции
curl http://localhost:8080/api/cancellation/instructions

# Поиск по названию
curl "http://localhost:8080/api/cancellation/instructions?search=spotify"

# Получить инструкцию для конкретного сервиса
curl http://localhost:8080/api/cancellation/instructions/Spotify

# Отменить подписку (требует токен)
curl -X POST \
  http://localhost:8080/api/subscriptions/1/cancel \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Frontend

**Страница `/cancel-subscriptions`:**
- Поиск по названию сервиса
- Аккордеон со списком сервисов
- Инструкции по отмене для каждого сервиса
- Ссылки на страницы отмены

**Функции JavaScript:**
- `loadServices()` - Загрузка всех сервисов из API
- `filterServices()` - Фильтрация по поисковому запросу
- `toggleService(id)` - Открытие/закрытие инструкции

### 4. Стили

**Файл `cancellation.css`:**
- Адаптивный дизайн
- Анимации открытия/закрытия
- Мобильная версия

## Предустановленные сервисы

Миграция добавляет инструкции для следующих популярных сервисов:

### Русские сервисы:
- Кинопоиск HD
- Яндекс Плюс
- ВК Музыка
- Облако Mail.ru
- Тинькофф Черный
- 1С

### Международные сервисы:
- Spotify
- Netflix
- YouTube Premium
- Telegram Premium
- Adobe
- Dropbox
- Google One
- iCloud
- LinkedIn Premium
- Notion
- ChatGPT Plus
- GitHub Pro
- 1Password
- Box
- Canva
- Zoom

## Добавление новых сервисов

Добавить новый сервис можно через SQL:

```sql
INSERT INTO cancellation_instructions
(service_name, instructions, note, url, category)
VALUES (
    'Название Сервиса',
    'Пошаговая инструкция по отмене подписки',
    'Дополнительные примечания (опционально)',
    'https://service.com/cancel',
    'Категория'
);
```

Или через API (если добавить соответствующий endpoint).

## Интеграция с Dashboard

Для интеграции с панелью пользователя:

1. Добавить кнопку "Отменить" на каждую подписку
2. При клике показывать модальное окно с инструкцией
3. Вызвать `POST /api/subscriptions/{id}/cancel` для отметки подписки как неактивной

```html
<button onclick="cancelSubscription(subscriptionId, serviceName)">
    Отменить
</button>
```

```javascript
async function cancelSubscription(subId, serviceName) {
    // 1. Получить инструкцию для сервиса
    const instruction = await fetch(`/api/cancellation/instructions/${serviceName}`);

    // 2. Показать модальное окно с инструкцией
    showCancellationModal(instruction);

    // 3. При подтверждении - отменить подписку
    await fetch(`/api/subscriptions/${subId}/cancel`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
}
```

## Тестирование

```bash
# 1. Запустить сервер
go run ./cmd/server/main.go

# 2. Открыть браузер
open http://localhost:8080/cancel-subscriptions

# 3. Протестировать поиск
# 4. Протестировать открытие/закрытие инструкций
# 5. Протестировать ссылки на страницы отмены
```

## Будущие улучшения

- [ ] Admin панель для добавления/редактирования сервисов
- [ ] Автоматическая отмена подписок через API сервисов (где доступно)
- [ ] Статистика популярных отменяемых подписок
- [ ] Уведомления об успешной отмене
- [ ] Интеграция с Telegram для получения инструкций через бота
