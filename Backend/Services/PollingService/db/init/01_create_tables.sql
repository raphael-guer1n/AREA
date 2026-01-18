CREATE TABLE IF NOT EXISTS polling_subscriptions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    action_id INTEGER NOT NULL UNIQUE,
    provider VARCHAR(64) NOT NULL,
    service VARCHAR(64) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    config JSONB NOT NULL DEFAULT '{}'::jsonb,
    interval_seconds INTEGER NOT NULL,
    last_item_id TEXT,
    last_polled_at TIMESTAMP,
    next_run_at TIMESTAMP,
    last_error TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_polling_subscriptions_user_id ON polling_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_polling_subscriptions_next_run ON polling_subscriptions(next_run_at);
CREATE INDEX IF NOT EXISTS idx_polling_subscriptions_active ON polling_subscriptions(active);
