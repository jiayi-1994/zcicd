-- Add missing clusters columns used by system model/service
ALTER TABLE clusters
ADD COLUMN IF NOT EXISTS kube_config_ref VARCHAR(256),
ADD COLUMN IF NOT EXISTS node_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS version VARCHAR(32);
