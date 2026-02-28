-- M4: Add parallel execution support to workflow stages
ALTER TABLE workflow_stages ADD COLUMN IF NOT EXISTS parallel BOOLEAN NOT NULL DEFAULT false;
