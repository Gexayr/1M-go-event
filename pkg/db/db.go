package db

import (
	"log"

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

	if err := database.AutoMigrate(&models.Event{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return database
}
