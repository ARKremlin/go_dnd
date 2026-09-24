CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users
(
    id                UUID PRIMARY KEY     DEFAULT uuid_generate_v4(),
    telegram_id       BIGINT,
    telegram_username TEXT,
    username          TEXT,
    password_hash     TEXT,
    role              TEXT        NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_role_check CHECK (role in ('player', 'dm')),
    CONSTRAINT users_role_consistency CHECK (
        CASE
            WHEN role = 'player' THEN
                telegram_id IS NOT NULL
                    AND password_hash IS NULL
                    AND username IS NULL
            WHEN role = 'dm' THEN
                telegram_id IS NULL
                    AND telegram_username IS NULL
                    AND password_hash IS NOT NULL
                    AND username IS NOT NULL
            ELSE FALSE
            END
        )
);

CREATE INDEX idx_users_telegram_id ON users (telegram_id) WHERE telegram_id IS NOT NULL;
CREATE INDEX idx_users_telegram_username ON users (telegram_username) WHERE telegram_username IS NOT NULL;
CREATE INDEX idx_users_username ON users (username) WHERE username IS NOT NULL;