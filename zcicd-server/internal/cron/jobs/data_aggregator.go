package jobs

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// DataAggregator computes dashboard statistics periodically.
type DataAggregator struct {
	db *gorm.DB
}

func NewDataAggregator(db *gorm.DB) *DataAggregator {
	return &DataAggregator{db: db}
}

func (d *DataAggregator) Run() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	d.db.Exec(`
		INSERT INTO daily_stats (stat_date, total_builds, successful_builds, failed_builds, total_deploys, successful_deploys, failed_deploys, updated_at)
		SELECT
			?::date,
			COALESCE((SELECT COUNT(*) FROM build_runs WHERE created_at >= ?), 0),
			COALESCE((SELECT COUNT(*) FROM build_runs WHERE created_at >= ? AND status = 'succeeded'), 0),
			COALESCE((SELECT COUNT(*) FROM build_runs WHERE created_at >= ? AND status = 'failed'), 0),
			COALESCE((SELECT COUNT(*) FROM deploy_histories WHERE created_at >= ?), 0),
			COALESCE((SELECT COUNT(*) FROM deploy_histories WHERE created_at >= ? AND status = 'succeeded'), 0),
			COALESCE((SELECT COUNT(*) FROM deploy_histories WHERE created_at >= ? AND status = 'failed'), 0),
			NOW()
		ON CONFLICT (stat_date) DO UPDATE SET
			total_builds = EXCLUDED.total_builds,
			successful_builds = EXCLUDED.successful_builds,
			failed_builds = EXCLUDED.failed_builds,
			total_deploys = EXCLUDED.total_deploys,
			successful_deploys = EXCLUDED.successful_deploys,
			failed_deploys = EXCLUDED.failed_deploys,
			updated_at = NOW()
	`, today, today, today, today, today, today, today)

	if d.db.Error != nil {
		log.Printf("aggregator: error: %v", d.db.Error)
	}
}
