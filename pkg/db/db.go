package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-event-registration/configs"
	"go-event-registration/internal/models"
)

// Init opens a DB connection and performs migrations.
func Init(cfg *configs.Config) *gorm.DB {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("failed to access database pool: %v", err)
	}
	// Conservative connection pool to prevent handshake overload
	sqlDB.SetMaxOpenConns(15)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := database.AutoMigrate(&models.Event{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return database
}
