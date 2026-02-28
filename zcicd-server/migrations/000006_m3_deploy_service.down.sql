-- ============================================================
-- M3: Rollback deploy service tables
-- ============================================================

-- Drop indexes
DROP INDEX IF EXISTS idx_env_variables_environment;
DROP INDEX IF EXISTS idx_approval_records_approver_status;
DROP INDEX IF EXISTS idx_approval_records_history;
DROP INDEX IF EXISTS idx_deploy_histories_status;
DROP INDEX IF EXISTS idx_deploy_histories_config_created;
DROP INDEX IF EXISTS idx_deploy_configs_argo_app;
DROP INDEX IF EXISTS idx_deploy_configs_environment;
DROP INDEX IF EXISTS idx_deploy_configs_service;
DROP INDEX IF EXISTS idx_deploy_configs_project;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS env_resource_quotas;
DROP TABLE IF EXISTS env_variables;
DROP TABLE IF EXISTS approval_records;
DROP TABLE IF EXISTS deploy_histories;
DROP TABLE IF EXISTS deploy_configs;
