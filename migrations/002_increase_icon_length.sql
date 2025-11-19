-- Increase icon column length to support longer URLs
ALTER TABLE subscriptions ALTER COLUMN icon TYPE VARCHAR(500);
