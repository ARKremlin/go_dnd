DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_telegram_id;
DROP INDEX IF EXISTS idx_users_telegram_username;

CREATE UNIQUE INDEX idx_users_username_unique
    ON users (username)
    WHERE username IS NOT NULL;

CREATE UNIQUE INDEX idx_users_telegram_id_unique
    ON users (telegram_id)
    WHERE telegram_id IS NOT NULL;

CREATE UNIQUE INDEX idx_users_telegram_username_unique
    ON users (telegram_username)
    WHERE telegram_username IS NOT NULL;