-- ============================================================
-- ZCI/CD Platform - Rollback Initial Schema
-- ============================================================

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS system_configs;
DROP TABLE IF EXISTS secrets;
DROP TABLE IF EXISTS audit_logs_2026;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS scan_results;
DROP TABLE IF EXISTS test_reports;
DROP TABLE IF EXISTS artifacts;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS stage_runs;
DROP TABLE IF EXISTS workflow_runs;
DROP TABLE IF EXISTS workflow_stages;
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS environments;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS clusters;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS env_type;

DROP EXTENSION IF EXISTS "uuid-ossp";
