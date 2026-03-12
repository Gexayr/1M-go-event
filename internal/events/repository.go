package events

import (
	"go-event-registration/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetEventByID(id uint) (*models.Event, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetEventByID(id uint) (*models.Event, error) {
	var event models.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}
