-- Debug script для проверки почему не работают уведомления

-- 1. Проверяем настройки Telegram
SELECT 'TELEGRAM SETTINGS:' as check_name;
SELECT user_id, chat_id, is_enabled, created_at
FROM telegram_settings;

-- 2. Проверяем подписки
SELECT '' as separator;
SELECT 'SUBSCRIPTIONS:' as check_name;
SELECT id, user_id, name, next_billing_date, is_active,
       (next_billing_date - CURRENT_DATE) as days_until_payment
FROM subscriptions;

-- 3. Проверяем настройки уведомлений
SELECT '' as separator;
SELECT 'NOTIFICATION SETTINGS:' as check_name;
SELECT n.subscription_id, s.name as subscription_name,
       n.notify_days_before, n.is_enabled
FROM notifications n
JOIN subscriptions s ON s.id = n.subscription_id;

-- 4. ГЛАВНАЯ ПРОВЕРКА: Какие уведомления должны отправиться
SELECT '' as separator;
SELECT 'NOTIFICATIONS THAT SHOULD BE SENT:' as check_name;
SELECT
    s.id as subscription_id,
    s.name as subscription_name,
    s.user_id,
    s.next_billing_date,
    n.notify_days_before,
    (s.next_billing_date - CURRENT_DATE) as days_until_payment,
    ts.chat_id,
    ts.is_enabled as telegram_enabled,
    CASE
        WHEN (s.next_billing_date - CURRENT_DATE) = n.notify_days_before THEN 'YES - SHOULD SEND!'
        ELSE 'NO - date mismatch'
    END as should_send
FROM subscriptions s
JOIN notifications n ON n.subscription_id = s.id
JOIN telegram_settings ts ON ts.user_id = s.user_id
WHERE s.is_active = true
  AND n.is_enabled = true
  AND ts.is_enabled = true;

-- 5. Проверяем историю уже отправленных уведомлений
SELECT '' as separator;
SELECT 'NOTIFICATION HISTORY:' as check_name;
SELECT nh.id, s.name as subscription_name,
       nh.sent_at, nh.notification_type
FROM notification_history nh
JOIN subscriptions s ON s.id = nh.subscription_id
ORDER BY nh.sent_at DESC
LIMIT 10;
