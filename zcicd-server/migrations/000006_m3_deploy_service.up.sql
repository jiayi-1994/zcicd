-- ============================================================
-- M3: Deploy service tables (Argo CD integration)
-- PostgreSQL 16
-- ============================================================

-- Deploy Configs (Argo CD Application mapping)
CREATE TABLE deploy_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id      UUID NOT NULL REFERENCES services(id),
    environment_id  UUID NOT NULL REFERENCES environments(id),
    name            VARCHAR(128) NOT NULL,
    deploy_type     VARCHAR(32) NOT NULL DEFAULT 'helm',  -- helm/kustomize/yaml
    repo_url        VARCHAR(512) NOT NULL,
    target_revision VARCHAR(128) DEFAULT 'main',
    chart_path      VARCHAR(256),
    values_override JSONB DEFAULT '{}',
    sync_policy     VARCHAR(32) DEFAULT 'manual',  -- manual/auto
    auto_sync       BOOLEAN DEFAULT FALSE,
    self_heal       BOOLEAN DEFAULT FALSE,
    prune           BOOLEAN DEFAULT FALSE,
    argo_app_name   VARCHAR(256),
    namespace       VARCHAR(128),
    status          VARCHAR(32) DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, service_id, environment_id)
);

-- Deploy Histories (deployment history records)
CREATE TABLE deploy_histories (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deploy_config_id  UUID NOT NULL REFERENCES deploy_configs(id) ON DELETE CASCADE,
    revision          VARCHAR(128),
    status            VARCHAR(32) NOT NULL DEFAULT 'pending',  -- pending/syncing/succeeded/failed/cancelled
    sync_status       VARCHAR(32),  -- Synced/OutOfSync/Unknown
    health_status     VARCHAR(32),  -- Healthy/Degraded/Progressing/Missing/Suspended/Unknown
    started_at        TIMESTAMPTZ,
    finished_at       TIMESTAMPTZ,
    duration          INTEGER,
    triggered_by      UUID REFERENCES users(id),
    rollback_from     UUID REFERENCES deploy_histories(id),
    gitops_commit     VARCHAR(64),
    diff_content      TEXT,
    error_message     TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Approval Records (deployment approval workflow)
CREATE TABLE approval_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deploy_history_id   UUID NOT NULL REFERENCES deploy_histories(id) ON DELETE CASCADE,
    environment_id      UUID NOT NULL REFERENCES environments(id),
    requested_by        UUID NOT NULL REFERENCES users(id),
    approver_id         UUID REFERENCES users(id),
    status              VARCHAR(32) NOT NULL DEFAULT 'pending',  -- pending/approved/rejected/expired
    comment             TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    decided_at          TIMESTAMPTZ
);

-- Environment Variables
CREATE TABLE env_variables (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment_id  UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    var_key         VARCHAR(256) NOT NULL,
    var_value       TEXT,
    is_secret       BOOLEAN DEFAULT FALSE,
    description     VARCHAR(512),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(environment_id, var_key)
);

-- Environment Resource Quotas
CREATE TABLE env_resource_quotas (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    environment_id  UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    cpu_request     VARCHAR(32),
    cpu_limit       VARCHAR(32),
    memory_request  VARCHAR(32),
    memory_limit    VARCHAR(32),
    pod_limit       INTEGER,
    storage_limit   VARCHAR(32),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(environment_id)
);

-- ============================================================
-- Indexes
-- ============================================================

-- deploy_configs
CREATE INDEX idx_deploy_configs_project ON deploy_configs(project_id);
CREATE INDEX idx_deploy_configs_service ON deploy_configs(service_id);
CREATE INDEX idx_deploy_configs_environment ON deploy_configs(environment_id);
CREATE INDEX idx_deploy_configs_argo_app ON deploy_configs(argo_app_name);

-- deploy_histories
CREATE INDEX idx_deploy_histories_config_created ON deploy_histories(deploy_config_id, created_at DESC);
CREATE INDEX idx_deploy_histories_status ON deploy_histories(status);

-- approval_records
CREATE INDEX idx_approval_records_history ON approval_records(deploy_history_id);
CREATE INDEX idx_approval_records_approver_status ON approval_records(approver_id, status);

-- env_variables
CREATE INDEX idx_env_variables_environment ON env_variables(environment_id);
