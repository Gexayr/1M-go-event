package alerts

import (
	"go-event-registration/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetAlerts() ([]models.Event, error)
	GetAlertStats() (int64, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetAlerts() ([]models.Event, error) {
	var events []models.Event
	err := r.db.Where("risk_score >= ?", 70).
		Order("timestamp DESC").
		Limit(50).
		Find(&events).Error
	return events, err
}

func (r *repository) GetAlertStats() (int64, int64, error) {
	var totalAlerts int64
	var criticalAlerts int64

	err := r.db.Model(&models.Event{}).Where("risk_score >= ?", 70).Count(&totalAlerts).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.Model(&models.Event{}).Where("risk_score >= ?", 85).Count(&criticalAlerts).Error
	if err != nil {
		return totalAlerts, 0, err
	}

	return totalAlerts, criticalAlerts, nil
}
