-- ============================================================
-- M2: Indexes for workflow + build tables
-- ============================================================

CREATE INDEX idx_stage_jobs_stage ON stage_jobs(stage_id, sort_order);
CREATE INDEX idx_build_configs_project ON build_configs(project_id);
CREATE INDEX idx_build_configs_service ON build_configs(service_id);
CREATE INDEX idx_build_runs_config ON build_runs(build_config_id);
CREATE INDEX idx_build_runs_status ON build_runs(status);
CREATE INDEX idx_build_templates_language ON build_templates(language);
