CREATE TABLE subscriptions (
                               id SERIAL PRIMARY KEY,
                               user_id BIGINT REFERENCES users(tg_id),
                               product_id INTEGER REFERENCES products(id),
                               target_price NUMERIC(12, 2),
                               is_triggered BOOLEAN DEFAULT FALSE,
                               UNIQUE(user_id, product_id)
);