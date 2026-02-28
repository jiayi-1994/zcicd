-- ============================================================
-- ZCI/CD Platform - Rollback Seed Data
-- ============================================================

DELETE FROM system_configs WHERE config_key IN (
    'site_name',
    'default_language',
    'session_timeout',
    'max_concurrent_builds',
    'build_log_retention_days',
    'audit_log_retention_days',
    'default_deploy_strategy',
    'enable_registration'
);

DELETE FROM user_roles WHERE user_id = '10000000-0000-0000-0000-000000000001';

DELETE FROM users WHERE id = '10000000-0000-0000-0000-000000000001';
