package database

import (
	"fmt"
	"time"

	"github.com/zcicd/zcicd-server/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	return db, nil
}

// DBPair holds a write (primary) and read (replica) DB connection.
type DBPair struct {
	Writer *gorm.DB
	Reader *gorm.DB
}

// NewPostgresWithReplicas creates a primary + replica DB pair for read/write split.
// Falls back to primary for reads if no replicas are configured.
func NewPostgresWithReplicas(cfg *config.Config) (*DBPair, error) {
	writer, err := NewPostgres(cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Database.Replicas) == 0 {
		return &DBPair{Writer: writer, Reader: writer}, nil
	}

	r := cfg.Database.Replicas[0]
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		r.Host, r.Port,
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)
	reader, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		// Fall back to primary if replica unavailable
		return &DBPair{Writer: writer, Reader: writer}, nil
	}

	sqlDB, _ := reader.DB()
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	return &DBPair{Writer: writer, Reader: reader}, nil
}
