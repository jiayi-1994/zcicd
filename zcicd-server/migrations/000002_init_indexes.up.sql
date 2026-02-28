-- ============================================================
-- ZCI/CD Platform - Indexes
-- ============================================================

-- Users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- User Roles
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);

-- Projects
CREATE INDEX idx_projects_identifier ON projects(identifier);
CREATE INDEX idx_projects_owner ON projects(owner_id);

-- Services
CREATE INDEX idx_services_project ON services(project_id);

-- Environments
CREATE INDEX idx_environments_project ON environments(project_id);

-- Workflows
CREATE INDEX idx_workflows_project ON workflows(project_id);

-- Workflow Runs
CREATE INDEX idx_workflow_runs_workflow ON workflow_runs(workflow_id);
CREATE INDEX idx_workflow_runs_status ON workflow_runs(status);

-- Stage Runs
CREATE INDEX idx_stage_runs_workflow_run ON stage_runs(workflow_run_id);

-- Deployments
CREATE INDEX idx_deployments_env ON deployments(environment_id);
CREATE INDEX idx_deployments_service ON deployments(service_id);

-- Artifacts
CREATE INDEX idx_artifacts_project ON artifacts(project_id);

-- Audit Logs
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_project ON audit_logs(project_id);

-- Secrets
CREATE INDEX idx_secrets_project ON secrets(project_id);

-- Notifications
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_unread ON notifications(user_id, is_read) WHERE is_read = false;
