CREATE TABLE IF NOT EXISTS webhook_subscriptions (
  id SERIAL PRIMARY KEY,
  hook_id VARCHAR(64) UNIQUE NOT NULL,
  provider_hook_id VARCHAR(128),
  user_id INTEGER NOT NULL,
  area_id INTEGER NOT NULL,
  provider VARCHAR(64) NOT NULL,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_hook_id ON webhook_subscriptions (hook_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_user_id ON webhook_subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_area_id ON webhook_subscriptions (area_id);
