-- Add Telegram authentication fields to users table
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS telegram_user_id BIGINT UNIQUE,
    ADD COLUMN IF NOT EXISTS telegram_username VARCHAR(100),
    ADD COLUMN IF NOT EXISTS telegram_first_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS telegram_last_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS telegram_photo_url TEXT,
    ADD COLUMN IF NOT EXISTS auth_type VARCHAR(20) DEFAULT 'email' CHECK (auth_type IN ('email', 'telegram')),
    ALTER COLUMN email DROP NOT NULL,
    ALTER COLUMN password_hash DROP NOT NULL;

-- Add constraint: email required for email auth, telegram_user_id required for telegram auth
ALTER TABLE users ADD CONSTRAINT users_auth_check
    CHECK (
        (auth_type = 'email' AND email IS NOT NULL AND password_hash IS NOT NULL) OR
        (auth_type = 'telegram' AND telegram_user_id IS NOT NULL)
    );

-- Create index for telegram_user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_telegram_id ON users(telegram_user_id);

-- Update existing users to have email auth type (for backward compatibility)
UPDATE users SET auth_type = 'email' WHERE auth_type IS NULL;
