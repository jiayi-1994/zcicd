-- ============================================================
-- M5: Quality & Artifact service tables
-- PostgreSQL 16
-- ============================================================

-- Test Configs
CREATE TABLE test_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    test_type       VARCHAR(32) NOT NULL DEFAULT 'unit',  -- unit/integration/e2e/performance
    framework       VARCHAR(64),          -- jest/pytest/go-test/jmeter
    command         TEXT,                 -- test execution command
    timeout         INTEGER DEFAULT 3600, -- seconds
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Test Runs
CREATE TABLE test_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_config_id  UUID NOT NULL REFERENCES test_configs(id) ON DELETE CASCADE,
    build_run_id    UUID,                 -- optional link to build
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',  -- pending/running/passed/failed/error
    total           INTEGER DEFAULT 0,
    passed          INTEGER DEFAULT 0,
    failed          INTEGER DEFAULT 0,
    skipped         INTEGER DEFAULT 0,
    coverage        NUMERIC(5,2),         -- code coverage percentage
    duration        INTEGER,              -- milliseconds
    report_url      VARCHAR(512),         -- MinIO path to full report
    error_message   TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Scan Configs (SonarQube integration)
CREATE TABLE scan_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name            VARCHAR(128) NOT NULL,
    scan_type       VARCHAR(32) NOT NULL DEFAULT 'sonar',  -- sonar/semgrep/trivy-fs
    sonar_project_key VARCHAR(256),
    config          JSONB DEFAULT '{}',   -- scanner-specific config
    enabled         BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Scan Runs
CREATE TABLE scan_runs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scan_config_id  UUID NOT NULL REFERENCES scan_configs(id) ON DELETE CASCADE,
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',  -- pending/running/completed/failed
    bugs            INTEGER DEFAULT 0,
    vulnerabilities INTEGER DEFAULT 0,
    code_smells     INTEGER DEFAULT 0,
    coverage        NUMERIC(5,2),
    duplications    NUMERIC(5,2),
    quality_rating  VARCHAR(1),           -- A/B/C/D/E
    gate_status     VARCHAR(16),          -- passed/failed
    report_url      VARCHAR(512),
    error_message   TEXT,
    started_at      TIMESTAMPTZ,
    finished_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Quality Gates
CREATE TABLE quality_gates (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id      UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    min_coverage    NUMERIC(5,2) DEFAULT 80.00,
    max_bugs        INTEGER DEFAULT 0,
    max_vulnerabilities INTEGER DEFAULT 0,
    max_code_smells INTEGER DEFAULT 50,
    max_duplications NUMERIC(5,2) DEFAULT 5.00,
    block_deploy    BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id)
);

-- Image Registries
CREATE TABLE image_registries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(128) NOT NULL,
    registry_type   VARCHAR(32) NOT NULL DEFAULT 'harbor',  -- harbor/dockerhub/acr/ecr
    endpoint        VARCHAR(512) NOT NULL,
    username        VARCHAR(128),
    password_enc    BYTEA,
    is_default      BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Image Scans (Trivy)
CREATE TABLE image_scans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    registry_id     UUID NOT NULL REFERENCES image_registries(id) ON DELETE CASCADE,
    image_name      VARCHAR(512) NOT NULL,
    tag             VARCHAR(256) NOT NULL,
    status          VARCHAR(32) NOT NULL DEFAULT 'pending',  -- pending/scanning/completed/failed
    critical        INTEGER DEFAULT 0,
    high            INTEGER DEFAULT 0,
    medium          INTEGER DEFAULT 0,
    low             INTEGER DEFAULT 0,
    report_url      VARCHAR(512),
    scanned_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Helm Charts
CREATE TABLE helm_charts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(256) NOT NULL,
    repo_url        VARCHAR(512) NOT NULL,
    latest_version  VARCHAR(64),
    description     TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_test_configs_project ON test_configs(project_id);
CREATE INDEX idx_test_runs_config ON test_runs(test_config_id, created_at DESC);
CREATE INDEX idx_test_runs_status ON test_runs(status);
CREATE INDEX idx_scan_configs_project ON scan_configs(project_id);
CREATE INDEX idx_scan_runs_config ON scan_runs(scan_config_id, created_at DESC);
CREATE INDEX idx_image_scans_registry ON image_scans(registry_id);
CREATE INDEX idx_image_scans_image ON image_scans(image_name, tag);
