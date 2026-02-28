-- ============================================================
-- M2: Workflow extensions + Build management tables
-- PostgreSQL 16
-- ============================================================

-- Stage Jobs (extends workflow_stages with individual jobs)
CREATE TABLE stage_jobs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stage_id        UUID NOT NULL REFERENCES workflow_stages(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    job_type        VARCHAR(32) NOT NULL,  -- build, test, deploy, custom, approval
    sort_order      INTEGER NOT NULL DEFAULT 0,
    config          JSONB NOT NULL DEFAULT '{}',
    timeout_sec     INTEGER DEFAULT 3600,
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Build Templates
CREATE TABLE build_templates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    language        VARCHAR(32) NOT NULL,  -- go, java, nodejs, python, etc.
    framework       VARCHAR(64),           -- gin, spring-boot, express, etc.
    description     TEXT,
    build_env       JSONB NOT NULL DEFAULT '{}',  -- base image, resources
    build_script    TEXT NOT NULL,          -- default build script
    dockerfile_tpl  TEXT,                   -- Dockerfile template
    tekton_task_tpl TEXT NOT NULL,          -- Tekton Task YAML template
    is_system       BOOLEAN DEFAULT FALSE, -- system preset vs user-created
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Build Configs (per service build configuration)
CREATE TABLE build_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    service_id      UUID NOT NULL REFERENCES services(id),
    name            VARCHAR(128) NOT NULL,
    template_id     UUID REFERENCES build_templates(id),
    repo_url        VARCHAR(512) NOT NULL,
    branch          VARCHAR(128) DEFAULT 'main',
    build_env       JSONB NOT NULL DEFAULT '{}',
    build_script    TEXT,
    dockerfile_path VARCHAR(256) DEFAULT 'Dockerfile',
    docker_context  VARCHAR(256) DEFAULT '.',
    image_repo      VARCHAR(256) NOT NULL,
    tag_strategy    VARCHAR(32) DEFAULT 'branch-commit',
    cache_enabled   BOOLEAN DEFAULT TRUE,
    variables       JSONB NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Build Runs (individual build execution records)
CREATE TABLE build_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    build_config_id UUID NOT NULL REFERENCES build_configs(id),
    run_number      INTEGER NOT NULL,
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',
    branch          VARCHAR(128),
    commit_sha      VARCHAR(64),
    commit_message  TEXT,
    image_tag       VARCHAR(256),
    image_digest    VARCHAR(128),
    tekton_ref      VARCHAR(256),
    log_path        VARCHAR(512),
    triggered_by    UUID REFERENCES users(id),
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    duration_sec    INTEGER,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Extend workflow_runs with tekton refs and input params
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS input_params JSONB NOT NULL DEFAULT '{}';
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS tekton_refs JSONB NOT NULL DEFAULT '{}';
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS stages_status JSONB NOT NULL DEFAULT '[]';
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS duration_sec INTEGER;
ALTER TABLE workflow_runs ADD COLUMN IF NOT EXISTS triggered_by UUID REFERENCES users(id);
