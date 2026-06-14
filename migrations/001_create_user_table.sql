CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username TEXT,
                       tg_id BIGINT UNIQUE,
                       premium_expires_at TIMESTAMP,
                       status TEXT DEFAULT 'free',
                       limit_prod INT DEFAULT 10
);