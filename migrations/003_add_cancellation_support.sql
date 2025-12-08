-- Add cancelled_at field to subscriptions
ALTER TABLE subscriptions ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMP;

-- Create cancellation_instructions table
CREATE TABLE IF NOT EXISTS cancellation_instructions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL UNIQUE,
    service_name_lower VARCHAR(255) NOT NULL UNIQUE, -- for case-insensitive search
    instructions TEXT NOT NULL,
    notes TEXT,
    cancellation_url VARCHAR(500),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for search
CREATE INDEX IF NOT EXISTS idx_cancellation_instructions_name ON cancellation_instructions(service_name_lower);

-- Insert popular services cancellation instructions
INSERT INTO cancellation_instructions (service_name, service_name_lower, instructions, notes, cancellation_url, is_verified) VALUES
('1Password', '1password', 'Войдите в свою учетную запись на 1Password.com → Выберите Billing (Оплата) в боковом меню → Выберите Cancel subscription (Отменить подписку)', 'Ваша учетная запись останется активной до конца текущего расчетного периода.', 'https://1password.com/account', true),
('Adobe', 'adobe', 'Войдите в Adobe Account → Перейдите в раздел Plans → Выберите Cancel plan → Подтвердите отмену', 'Для подписки через App Store / Google Play отмена только в магазине.', 'https://account.adobe.com/plans', true),
('Box', 'box', 'Войдите в Box → Settings → Account → Billing → Cancel Subscription', NULL, 'https://app.box.com/settings/account', true),
('Canva', 'canva', 'Войдите в Canva → Account settings → Billing → Cancel subscription', NULL, 'https://www.canva.com/account', true),
('Облако Mail.ru', 'облако mail.ru', 'Личный кабинет Облака → Настройки → Тариф → Отменить подписку', NULL, 'https://cloud.mail.ru', true),
('Netflix', 'netflix', 'Войдите в Netflix → Account → Membership & Billing → Cancel Membership', 'Ваша подписка будет активна до конца текущего периода оплаты.', 'https://www.netflix.com/YourAccount', true),
('Spotify', 'spotify', 'Войдите в Spotify → Account → Subscription → Cancel Subscription', 'Для подписки через App Store / Google Play отмена только в магазине.', 'https://www.spotify.com/account/subscription', true),
('YouTube Premium', 'youtube premium', 'Войдите в YouTube → Settings → Membership → Cancel membership', 'Для подписки через App Store / Google Play отмена только в магазине.', 'https://www.youtube.com/paid_memberships', true),
('Kinopoisk', 'kinopoisk', 'Личный кабинет Кинопоиска → Подписка → Отменить подписку', NULL, 'https://www.kinopoisk.ru/premium', true),
('Яндекс Плюс', 'яндекс плюс', 'Личный кабинет Плюса → Отменить подписку → Подтвердить до экрана «Вы отменили подписку»', 'Для подписки через App Store / Google Play отмена только в магазине.', 'https://plus.yandex.ru', true),
('Apple Music', 'apple music', 'На iPhone/iPad: Settings → [ваше имя] → Subscriptions → Apple Music → Cancel Subscription. На Android: Откройте приложение Apple Music → Account → Manage Subscription → Cancel Subscription', 'Для подписки через App Store / Google Play отмена только в магазине.', NULL, true),
('ChatGPT Plus', 'chatgpt plus', 'Войдите в ChatGPT → Settings → Billing → Manage subscription → Cancel plan', NULL, 'https://chat.openai.com/settings/billing', true),
('GitHub Copilot', 'github copilot', 'Войдите в GitHub → Settings → Billing → GitHub Copilot → Cancel subscription', NULL, 'https://github.com/settings/billing', true),
('Yandex 360', 'yandex 360', 'Личный кабинет 360 → Тарифы → Отменить подписку', NULL, 'https://360.yandex.ru', true),
('Notion', 'notion', 'Войдите в Notion → Settings & Members → Billing → Cancel subscription', NULL, 'https://www.notion.so/settings/billing', true),
('1C', '1c', 'Личный кабинет 1С → Подписки → Отменить подписку', NULL, 'https://online.1c.ru', true)
ON CONFLICT (service_name) DO NOTHING;

-- Create trigger for updated_at
CREATE TRIGGER update_cancellation_instructions_updated_at BEFORE UPDATE ON cancellation_instructions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

