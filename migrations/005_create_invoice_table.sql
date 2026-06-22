CREATE TABLE invoice(
                        invoice_id TEXT PRIMARY KEY,
                        bank_name TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        amount NUMERIC,
                        user_id BIGINT REFERENCES users(id)
);