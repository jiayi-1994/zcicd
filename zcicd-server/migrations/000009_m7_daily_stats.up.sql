CREATE TABLE IF NOT EXISTS daily_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stat_date DATE NOT NULL UNIQUE,
    total_builds INT DEFAULT 0,
    successful_builds INT DEFAULT 0,
    failed_builds INT DEFAULT 0,
    total_deploys INT DEFAULT 0,
    successful_deploys INT DEFAULT 0,
    failed_deploys INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
