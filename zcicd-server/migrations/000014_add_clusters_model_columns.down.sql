-- Roll back clusters model column additions
ALTER TABLE clusters
DROP COLUMN IF EXISTS version,
DROP COLUMN IF EXISTS node_count,
DROP COLUMN IF EXISTS kube_config_ref;
