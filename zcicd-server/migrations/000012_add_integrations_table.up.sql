-- Add missing integrations table for system module
CREATE TABLE IF NOT EXISTS integrations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    type            VARCHAR(32) NOT NULL,
    provider        VARCHAR(32) NOT NULL,
    config_enc      BYTEA NOT NULL,
    status          VARCHAR(16) DEFAULT 'active',
    last_check_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integrations_type_provider ON integrations(type, provider);
CREATE INDEX IF NOT EXISTS idx_integrations_status ON integrations(status);
