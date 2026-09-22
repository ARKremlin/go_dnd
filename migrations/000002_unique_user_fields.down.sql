DROP INDEX IF EXISTS idx_users_username_unique;
DROP INDEX IF EXISTS idx_users_telegram_id_unique;
DROP INDEX IF EXISTS idx_users_telegram_username_unique;

CREATE INDEX idx_users_username
    ON users (username)
    WHERE username IS NOT NULL;

CREATE INDEX idx_users_telegram_id
    ON users (telegram_id)
    WHERE telegram_id IS NOT NULL;

CREATE INDEX idx_users_telegram_username
    ON users (telegram_username)
    WHERE telegram_username IS NOT NULL;