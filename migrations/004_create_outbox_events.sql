CREATE TABLE outbox_events(
                              id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              payload JSONB,
                              status TEXT DEFAULT 'pending',
                              created_at TIMESTAMP DEFAULT NOW()


);