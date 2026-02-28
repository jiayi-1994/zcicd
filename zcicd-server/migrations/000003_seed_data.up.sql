-- ============================================================
-- ZCI/CD Platform - Seed Data
-- ============================================================

-- Insert admin user
-- Password: admin123 (bcrypt placeholder)
INSERT INTO users (id, username, email, password_hash, display_name, status, created_at, updated_at)
VALUES (
    '10000000-0000-0000-0000-000000000001',
    'admin',
    'admin@zcicd.local',
    '$2a$10$placeholder.hash.for.admin123.replace.in.production',
    'System Administrator',
    'active',
    NOW(),
    NOW()
);

-- Assign admin role
INSERT INTO user_roles (user_id, role, scope_type, scope_id)
VALUES (
    '10000000-0000-0000-0000-000000000001',
    'admin',
    'system',
    '00000000-0000-0000-0000-000000000000'
);

-- Default system configs
INSERT INTO system_configs (config_key, config_value, description, updated_at) VALUES
('site_name', '"ZCI/CD Platform"', 'Platform display name', NOW()),
('default_language', '"zh-CN"', 'Default UI language', NOW()),
('session_timeout', '3600', 'Session timeout in seconds', NOW()),
('max_concurrent_builds', '10', 'Maximum concurrent build runs', NOW()),
('build_log_retention_days', '30', 'Build log retention period in days', NOW()),
('audit_log_retention_days', '365', 'Audit log retention period in days', NOW()),
('default_deploy_strategy', '"rolling"', 'Default deployment strategy', NOW()),
('enable_registration', 'false', 'Whether to allow self-registration', NOW());
