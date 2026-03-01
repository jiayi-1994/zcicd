-- Rollback display_name column addition on clusters
ALTER TABLE clusters
DROP COLUMN IF EXISTS display_name;
