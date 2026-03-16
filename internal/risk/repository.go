package risk

import (
	"go-event-registration/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	GetClientsStats() ([]ClientStats, error)
	GetClientProfile(clientID string) (*ClientProfileResponse, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetClientsStats() ([]ClientStats, error) {
	var stats []ClientStats
	err := r.db.Model(&models.Event{}).
		Select("client_id, COUNT(*) as total_events, SUM(CASE WHEN risk_score >= 70 THEN 1 ELSE 0 END) as high_risk_events, AVG(risk_score) as avg_risk").
		Group("client_id").
		Order("avg_risk DESC").
		Scan(&stats).Error
	return stats, err
}

func (r *repository) GetClientProfile(clientID string) (*ClientProfileResponse, error) {
	var profile ClientProfileResponse
	highRiskThreshold := 75
	err := r.db.Model(&models.Event{}).
		Select(`
           COUNT(*) as total_events, 
           AVG(risk_score) as avg_risk,
           COUNT(CASE WHEN risk_score >= ? THEN 1 END) as high_risk_events
       `, highRiskThreshold).
		Where("client_id = ?", clientID).
		Scan(&profile).Error
	if err != nil {
		return nil, err
	}

	var recentEvents []models.Event
	err = r.db.Where("client_id = ?", clientID).
		Order("timestamp DESC").
		Limit(10).
		Find(&recentEvents).Error
	if err != nil {
		return nil, err
	}

	profile.RecentEvents = recentEvents
	profile.ClientID = clientID
	return &profile, nil
}
