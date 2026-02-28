-- ============================================================
-- M6: System service tables (notifications, dashboard, clusters, integrations, audit)
-- PostgreSQL 16
-- ============================================================

-- Notification Channels
CREATE TABLE notify_channels (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    channel_type    VARCHAR(32) NOT NULL,  -- dingtalk/wechat/slack/email/webhook
    config          JSONB NOT NULL DEFAULT '{}',
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notification Rules
CREATE TABLE notify_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    event_type      VARCHAR(64) NOT NULL,  -- build.completed/deploy.succeeded/approval.pending
    channel_id      UUID NOT NULL REFERENCES notify_channels(id) ON DELETE CASCADE,
    project_id      UUID REFERENCES projects(id) ON DELETE CASCADE,
    severity        VARCHAR(16) DEFAULT 'all',  -- all/critical/warning/info
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Notification History
CREATE TABLE notify_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    channel_id      UUID NOT NULL REFERENCES notify_channels(id),
    event_type      VARCHAR(64) NOT NULL,
    title           VARCHAR(256),
    content         TEXT,
    status          VARCHAR(16) NOT NULL DEFAULT 'pending',  -- pending/sent/failed
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_notify_rules_event ON notify_rules(event_type);
CREATE INDEX idx_notify_rules_channel ON notify_rules(channel_id);
CREATE INDEX idx_notify_history_channel ON notify_history(channel_id, created_at DESC);
CREATE INDEX idx_notify_history_status ON notify_history(status);
