-- ============================================================
-- M2: Rollback workflow extensions + build management tables
-- ============================================================

ALTER TABLE workflow_runs DROP COLUMN IF EXISTS input_params;
ALTER TABLE workflow_runs DROP COLUMN IF EXISTS tekton_refs;
ALTER TABLE workflow_runs DROP COLUMN IF EXISTS stages_status;
ALTER TABLE workflow_runs DROP COLUMN IF EXISTS duration_sec;
ALTER TABLE workflow_runs DROP COLUMN IF EXISTS triggered_by;

DROP TABLE IF EXISTS build_runs;
DROP TABLE IF EXISTS build_configs;
DROP TABLE IF EXISTS build_templates;
DROP TABLE IF EXISTS stage_jobs;
