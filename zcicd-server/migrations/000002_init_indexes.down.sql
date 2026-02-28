-- ============================================================
-- ZCI/CD Platform - Rollback Indexes
-- ============================================================

DROP INDEX IF EXISTS idx_notifications_unread;
DROP INDEX IF EXISTS idx_notifications_user;
DROP INDEX IF EXISTS idx_secrets_project;
DROP INDEX IF EXISTS idx_audit_logs_project;
DROP INDEX IF EXISTS idx_audit_logs_created;
DROP INDEX IF EXISTS idx_audit_logs_user;
DROP INDEX IF EXISTS idx_artifacts_project;
DROP INDEX IF EXISTS idx_deployments_service;
DROP INDEX IF EXISTS idx_deployments_env;
DROP INDEX IF EXISTS idx_stage_runs_workflow_run;
DROP INDEX IF EXISTS idx_workflow_runs_status;
DROP INDEX IF EXISTS idx_workflow_runs_workflow;
DROP INDEX IF EXISTS idx_workflows_project;
DROP INDEX IF EXISTS idx_environments_project;
DROP INDEX IF EXISTS idx_services_project;
DROP INDEX IF EXISTS idx_projects_owner;
DROP INDEX IF EXISTS idx_projects_identifier;
DROP INDEX IF EXISTS idx_user_roles_user_id;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
