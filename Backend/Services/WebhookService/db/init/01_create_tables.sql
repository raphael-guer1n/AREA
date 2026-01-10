CREATE TABLE IF NOT EXISTS webhook_subscriptions (
  id SERIAL PRIMARY KEY,
  hook_id VARCHAR(64) UNIQUE NOT NULL,
  provider_hook_id VARCHAR(128),
  user_id INTEGER NOT NULL,
  action_id INTEGER NOT NULL,
  provider VARCHAR(64) NOT NULL,
  service VARCHAR(64) NOT NULL,
  auth_token TEXT,
  active BOOLEAN NOT NULL DEFAULT TRUE,
  config JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_hook_id ON webhook_subscriptions (hook_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_user_id ON webhook_subscriptions (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_webhook_subscriptions_action_id ON webhook_subscriptions (action_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_service ON webhook_subscriptions (service);
