-- M6: Rollback system service tables
DROP TABLE IF EXISTS notify_history CASCADE;
DROP TABLE IF EXISTS notify_rules CASCADE;
DROP TABLE IF EXISTS notify_channels CASCADE;
