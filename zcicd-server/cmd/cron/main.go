package main

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/zcicd/zcicd-server/internal/cron/jobs"
	"github.com/zcicd/zcicd-server/pkg/config"
	"github.com/zcicd/zcicd-server/pkg/database"
	"github.com/zcicd/zcicd-server/pkg/logger"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	logger.Init(cfg)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	c := cron.New(cron.WithSeconds())

	cleaner := jobs.NewResourceCleaner(db, 7*24*time.Hour)
	aggregator := jobs.NewDataAggregator(db)

	// Clean old build/deploy history every day at 2:00 AM
	c.AddFunc("0 0 2 * * *", cleaner.Run)
	// Aggregate dashboard stats every 10 minutes
	c.AddFunc("0 */10 * * * *", aggregator.Run)

	c.Start()
	log.Println("cron-service started")

	select {}
}
