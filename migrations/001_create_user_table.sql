CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       username TEXT,
                       tg_id BIGINT UNIQUE

);