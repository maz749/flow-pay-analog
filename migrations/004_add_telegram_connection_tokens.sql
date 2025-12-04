-- Add table for temporary connection tokens
CREATE TABLE IF NOT EXISTS pending_connections (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(100) UNIQUE NOT NULL,
    expires_at TIMESTAMP DEFAULT (NOW() + INTERVAL '1 hour'),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Add index for faster token lookups
CREATE INDEX idx_pending_connections_token ON pending_connections(token);
CREATE INDEX idx_pending_connections_expires ON pending_connections(expires_at);

-- Add index to telegram_settings for faster lookups
CREATE INDEX IF NOT EXISTS idx_telegram_settings_user_id ON telegram_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_telegram_settings_chat_id ON telegram_settings(chat_id);

-- Add notification_logs table for tracking sent notifications
CREATE TABLE IF NOT EXISTS notification_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    subscription_id INTEGER REFERENCES subscriptions(id) ON DELETE SET NULL,
    chat_id BIGINT,
    message_text TEXT,
    status VARCHAR(20) NOT NULL, -- 'sent', 'failed', 'pending'
    error_message TEXT,
    sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notification_logs_user_id ON notification_logs(user_id);
CREATE INDEX idx_notification_logs_sent_at ON notification_logs(sent_at);
CREATE INDEX idx_notification_logs_status ON notification_logs(status);

-- Clean up expired tokens automatically (PostgreSQL function)
CREATE OR REPLACE FUNCTION cleanup_expired_tokens()
RETURNS void AS $$
BEGIN
    DELETE FROM pending_connections WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Optional: Create a scheduled job to clean expired tokens
-- This requires pg_cron extension (install separately if needed)
-- SELECT cron.schedule('cleanup-tokens', '0 * * * *', 'SELECT cleanup_expired_tokens()');
