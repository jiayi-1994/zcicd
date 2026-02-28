package jobs

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// ResourceCleaner removes old build/deploy history records.
type ResourceCleaner struct {
	db        *gorm.DB
	retention time.Duration
}

func NewResourceCleaner(db *gorm.DB, retention time.Duration) *ResourceCleaner {
	return &ResourceCleaner{db: db, retention: retention}
}

func (r *ResourceCleaner) Run() {
	cutoff := time.Now().Add(-r.retention)

	res := r.db.Exec("DELETE FROM build_runs WHERE status IN ('succeeded','failed') AND finished_at < ?", cutoff)
	if res.Error != nil {
		log.Printf("cleaner: build_runs error: %v", res.Error)
	} else if res.RowsAffected > 0 {
		log.Printf("cleaner: removed %d old build_runs", res.RowsAffected)
	}

	res = r.db.Exec("DELETE FROM deploy_histories WHERE status IN ('succeeded','failed') AND finished_at < ?", cutoff)
	if res.Error != nil {
		log.Printf("cleaner: deploy_histories error: %v", res.Error)
	} else if res.RowsAffected > 0 {
		log.Printf("cleaner: removed %d old deploy_histories", res.RowsAffected)
	}

	res = r.db.Exec("DELETE FROM scan_runs WHERE status IN ('completed','failed') AND created_at < ?", cutoff)
	if res.Error != nil {
		log.Printf("cleaner: scan_runs error: %v", res.Error)
	} else if res.RowsAffected > 0 {
		log.Printf("cleaner: removed %d old scan_runs", res.RowsAffected)
	}
}
