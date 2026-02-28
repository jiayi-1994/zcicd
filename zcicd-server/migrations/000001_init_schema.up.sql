-- ============================================================
-- ZCI/CD Platform - Initial Schema Migration
-- PostgreSQL 16
-- ============================================================

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================
-- Custom ENUM types
-- ============================================================

CREATE TYPE env_type AS ENUM ('dev', 'testing', 'staging', 'production');

-- ============================================================
-- Users & Roles
-- ============================================================

CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(64) UNIQUE NOT NULL,
    email           VARCHAR(255) UNIQUE NOT NULL,
    password_hash   VARCHAR(255),
    display_name    VARCHAR(128),
    avatar_url      VARCHAR(512),
    status          VARCHAR(16) NOT NULL DEFAULT 'active',
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_roles (
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role            VARCHAR(64) NOT NULL,
    scope_type      VARCHAR(32) NOT NULL DEFAULT 'system',
    scope_id        UUID NOT NULL DEFAULT '00000000-0000-0000-0000-000000000000',
    PRIMARY KEY (user_id, role, scope_type, scope_id)
);

-- ============================================================
-- Projects
-- ============================================================

CREATE TABLE projects (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    identifier      VARCHAR(128) UNIQUE NOT NULL,
    description     TEXT,
    owner_id        UUID NOT NULL REFERENCES users(id),
    repo_url        VARCHAR(512),
    default_branch  VARCHAR(128) DEFAULT 'main',
    visibility      VARCHAR(16) NOT NULL DEFAULT 'private',
    status          VARCHAR(16) NOT NULL DEFAULT 'active',
    settings        JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- ============================================================
-- Clusters (must be before environments)
-- ============================================================

CREATE TABLE clusters (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(128) UNIQUE NOT NULL,
    description         TEXT,
    kubeconfig_encrypted BYTEA,
    api_server_url      VARCHAR(512),
    status              VARCHAR(16) NOT NULL DEFAULT 'connected',
    provider            VARCHAR(32),
    region              VARCHAR(64),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Services
-- ============================================================

CREATE TABLE services (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    service_type    VARCHAR(32),
    language        VARCHAR(32),
    repo_url        VARCHAR(512),
    branch          VARCHAR(128) DEFAULT 'main',
    dockerfile_path VARCHAR(256) DEFAULT 'Dockerfile',
    build_context   VARCHAR(256) DEFAULT '.',
    deploy_type     VARCHAR(32) NOT NULL DEFAULT 'yaml',
    helm_chart_path VARCHAR(256),
    helm_values     JSONB,
    k8s_manifests   TEXT,
    health_check_path VARCHAR(256),
    ports           JSONB,
    env_vars        JSONB,
    resources       JSONB,
    status          VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Environments
-- ============================================================

CREATE TABLE environments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    env_type        env_type NOT NULL DEFAULT 'dev',
    namespace       VARCHAR(128),
    cluster_id      UUID REFERENCES clusters(id),
    is_production   BOOLEAN NOT NULL DEFAULT FALSE,
    auto_deploy     BOOLEAN NOT NULL DEFAULT FALSE,
    deploy_strategy JSONB,
    global_env_vars JSONB,
    status          VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Workflows
-- ============================================================

CREATE TABLE workflows (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    description     TEXT,
    trigger_type    VARCHAR(32),
    trigger_config  JSONB,
    enabled         BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE workflow_stages (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    stage_type      VARCHAR(32),
    sort_order      INTEGER NOT NULL DEFAULT 0,
    config          JSONB,
    timeout         INTEGER DEFAULT 3600,
    enabled         BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE workflow_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_id     UUID NOT NULL REFERENCES workflows(id),
    run_number      INTEGER NOT NULL,
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    trigger_type    VARCHAR(32) NOT NULL,
    trigger_info    JSONB,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration        INTEGER,
    error_message   TEXT,
    created_by      UUID REFERENCES users(id)
);

CREATE TABLE stage_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_run_id UUID NOT NULL REFERENCES workflow_runs(id) ON DELETE CASCADE,
    stage_id        UUID NOT NULL REFERENCES workflow_stages(id),
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    log_url         VARCHAR(512),
    error_message   TEXT
);

-- ============================================================
-- Deployments
-- ============================================================

CREATE TABLE deployments (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment_id  UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    service_id      UUID NOT NULL REFERENCES services(id),
    version         VARCHAR(128),
    image           VARCHAR(512),
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    deploy_type     VARCHAR(32),
    strategy        JSONB,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    rollback_from   UUID,
    created_by      UUID REFERENCES users(id),
    gitops_commit   VARCHAR(64)
);

-- ============================================================
-- Artifacts
-- ============================================================

CREATE TABLE artifacts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id      UUID REFERENCES services(id),
    workflow_run_id UUID REFERENCES workflow_runs(id),
    artifact_type   VARCHAR(32) NOT NULL,
    name            VARCHAR(256) NOT NULL,
    version         VARCHAR(128),
    storage_path    VARCHAR(512),
    size            BIGINT,
    checksum        VARCHAR(128),
    metadata        JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ
);

-- ============================================================
-- Test Reports
-- ============================================================

CREATE TABLE test_reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_run_id UUID NOT NULL REFERENCES workflow_runs(id),
    service_id      UUID REFERENCES services(id),
    test_type       VARCHAR(32),
    framework       VARCHAR(64),
    total           INTEGER NOT NULL DEFAULT 0,
    passed          INTEGER NOT NULL DEFAULT 0,
    failed          INTEGER NOT NULL DEFAULT 0,
    skipped         INTEGER NOT NULL DEFAULT 0,
    coverage        NUMERIC(5,2),
    report_url      VARCHAR(512),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Scan Results
-- ============================================================

CREATE TABLE scan_results (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workflow_run_id UUID NOT NULL REFERENCES workflow_runs(id),
    service_id      UUID REFERENCES services(id),
    scan_type       VARCHAR(32),
    scanner         VARCHAR(64),
    critical        INTEGER NOT NULL DEFAULT 0,
    high            INTEGER NOT NULL DEFAULT 0,
    medium          INTEGER NOT NULL DEFAULT 0,
    low             INTEGER NOT NULL DEFAULT 0,
    report_url      VARCHAR(512),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Audit Logs (Partitioned by created_at)
-- ============================================================

CREATE TABLE audit_logs (
    id              UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id         UUID,
    username        VARCHAR(64),
    action          VARCHAR(64) NOT NULL,
    resource_type   VARCHAR(64) NOT NULL,
    resource_id     UUID,
    project_id      UUID,
    detail          JSONB,
    ip_address      INET,
    user_agent      VARCHAR(512),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create initial partition for 2026
CREATE TABLE audit_logs_2026 PARTITION OF audit_logs
    FOR VALUES FROM ('2026-01-01') TO ('2027-01-01');

-- ============================================================
-- Secrets
-- ============================================================

CREATE TABLE secrets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    encrypted_value BYTEA NOT NULL,
    scope           VARCHAR(32) NOT NULL DEFAULT 'project',
    key_version     INTEGER NOT NULL DEFAULT 1,
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- System Configs
-- ============================================================

CREATE TABLE system_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_key      VARCHAR(128) UNIQUE NOT NULL,
    config_value    JSONB NOT NULL DEFAULT '{}',
    description     TEXT,
    updated_by      UUID REFERENCES users(id),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================
-- Notifications
-- ============================================================

CREATE TABLE notifications (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title                   VARCHAR(256) NOT NULL,
    content                 TEXT,
    notification_type       VARCHAR(32),
    is_read                 BOOLEAN NOT NULL DEFAULT FALSE,
    related_resource_type   VARCHAR(64),
    related_resource_id     UUID,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
