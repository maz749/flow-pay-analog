-- Create cancellation_instructions table
CREATE TABLE IF NOT EXISTS cancellation_instructions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) UNIQUE NOT NULL,
    instructions TEXT NOT NULL,
    note TEXT,
    url VARCHAR(500),
    category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for faster search
CREATE INDEX idx_cancellation_service_name ON cancellation_instructions(service_name);

-- Insert popular Russian services
INSERT INTO cancellation_instructions (service_name, instructions, note, url, category) VALUES
('1Password', 'Войдите в свою учетную запись на 1Password.com → Выберите Billing (Оплата) в боковом меню → Выберите Cancel subscription (Отменить подписку)', 'Ваша учетная запись останется активной до конца текущего расчетного периода.', 'https://my.1password.com/billing', 'Безопасность'),
('Adobe', 'Войдите на сайт Adobe → Управление планом → Отменить план → Выберите причину отмены → Подтвердите', 'Возможна комиссия за досрочную отмену, если подписка оформлена менее года назад.', 'https://account.adobe.com/plans', 'Дизайн'),
('Box', 'Войдите в Box → Настройки аккаунта → Выберите вкладку Аккаунт → Выберите Отменить аккаунт', 'Данные будут удалены через 7 дней после отмены.', 'https://account.box.com/', 'Облако'),
('Canva', 'Настройки → Аккаунт → Управление подпиской → Отменить подписку', 'Доступ к Pro-функциям сохранится до конца оплаченного периода.', 'https://www.canva.com/settings/billing', 'Дизайн'),
('ChatGPT Plus', 'Настройки → Моя подписка → Управление подпиской → Отменить план', 'Доступ к ChatGPT Plus сохранится до конца текущего периода.', 'https://chat.openai.com/', 'AI'),
('Dropbox', 'Настройки → Управление планом → Отменить подписку', 'После отмены вернётесь на бесплатный план с 2 ГБ хранилища.', 'https://www.dropbox.com/account/plan', 'Облако'),
('GitHub Pro', 'Настройки → Биллинг и планы → Изменить план → Downgrade to Free', 'Приватные репозитории станут публичными или будут удалены.', 'https://github.com/settings/billing', 'Разработка'),
('Google One', 'Google One приложение → Настройки → Отменить членство', 'Хранилище вернётся к 15 ГБ бесплатно.', 'https://one.google.com/storage', 'Облако'),
('iCloud', 'Настройки → Ваше имя → iCloud → Управлять хранилищем → Изменить план хранилища → Downgrade Options', 'Для подписки через App Store/Google Play отмена только в магазине.', '', 'Облако'),
('Кинопоиск HD', 'Личный кабинет Кинопоиска → Управление подпиской → Отменить подписку', 'Доступ сохранится до конца оплаченного периода.', 'https://hd.kinopoisk.ru/settings/subscription', 'Стриминг'),
('LinkedIn Premium', 'Настройки → Управление подпиской → Отменить подписку', 'Премиум-функции будут доступны до конца периода.', 'https://www.linkedin.com/psettings/subscriptions', 'Социальные'),
('Notion', 'Настройки → Планы и биллинг → Управление → Отменить план', 'Вернётесь на Free план с ограничениями.', 'https://www.notion.so/settings', 'Продуктивность'),
('Netflix', 'Аккаунт → Членство и биллинг → Отменить членство', 'Доступ сохранится до конца текущего периода.', 'https://www.netflix.com/YourAccount', 'Стриминг'),
('Облако Mail.ru', 'Мой профиль → Тарифы → Отключить подписку', 'Вернётесь на бесплатный тариф с 8 ГБ.', 'https://cloud.mail.ru/', 'Облако'),
('Spotify', 'Аккаунт → Управление подпиской → Отменить подписку', 'Премиум сохранится до конца оплаченного месяца.', 'https://www.spotify.com/account/subscription/', 'Музыка'),
('Telegram Premium', 'Настройки Telegram → Telegram Premium → Управлять подпиской', 'Для подписки через App Store/Google Play отмена только в магазине.', '', 'Мессенджеры'),
('Tinkoff Черный', 'Приложение Тинькофф → Профиль → Подписка → Отменить', 'Обслуживание прекратится с начала следующего месяца.', 'https://www.tinkoff.ru/', 'Финансы'),
('ВК Музыка', 'Настройки → Подписки → VK Музыка → Отключить автопродление', 'Доступ сохранится до конца периода.', 'https://vk.com/settings?act=payments', 'Музыка'),
('Яндекс Плюс', 'Личный кабинет Плюса → Отменить подписку → Подтвердить до экрана «Вы отменили подписку»', 'Для подписки через App Store / Google Play отмена только в магазине', 'https://plus.yandex.ru/', 'Подписки'),
('YouTube Premium', 'YouTube → Платные подписки → Управлять подпиской → Отменить подписку', 'Для подписки через App Store/Google Play отмена только в магазине.', 'https://www.youtube.com/paid_memberships', 'Стриминг'),
('Zoom', 'Управление аккаунтом → Биллинг → Отменить подписку', 'Вернётесь на бесплатный план с лимитом 40 минут для встреч.', 'https://zoom.us/billing', 'Видеосвязь');

-- Create trigger for updated_at
CREATE TRIGGER update_cancellation_instructions_updated_at BEFORE UPDATE ON cancellation_instructions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
