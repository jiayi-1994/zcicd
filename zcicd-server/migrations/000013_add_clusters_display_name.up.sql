-- Add missing display_name column to clusters table for API/model compatibility
ALTER TABLE clusters
ADD COLUMN IF NOT EXISTS display_name VARCHAR(256);
