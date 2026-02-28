-- ============================================================
-- M2: Drop indexes for workflow + build tables
-- ============================================================

DROP INDEX IF EXISTS idx_build_templates_language;
DROP INDEX IF EXISTS idx_build_runs_status;
DROP INDEX IF EXISTS idx_build_runs_config;
DROP INDEX IF EXISTS idx_build_configs_service;
DROP INDEX IF EXISTS idx_build_configs_project;
DROP INDEX IF EXISTS idx_stage_jobs_stage;
