-- Create cron_actions table
CREATE TABLE IF NOT EXISTS cron_actions (
    action_id INTEGER PRIMARY KEY,
    active BOOLEAN NOT NULL DEFAULT true,
    type TEXT NOT NULL,
    provider TEXT,
    service TEXT NOT NULL,
    title TEXT NOT NULL,
    input JSONB NOT NULL,
    cron_job_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on active actions for faster queries
CREATE INDEX idx_cron_actions_active ON cron_actions(active);

-- Create index on created_at for sorting
CREATE INDEX idx_cron_actions_created_at ON cron_actions(created_at DESC);
