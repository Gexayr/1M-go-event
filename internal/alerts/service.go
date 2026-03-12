package alerts

import (
	"go-event-registration/internal/models"
)

type Service interface {
	GetAlerts() ([]AlertResponse, error)
	GetAlertStats() (*AlertStatsResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAlerts() ([]AlertResponse, error) {
	events, err := s.repo.GetAlerts()
	if err != nil {
		return nil, err
	}

	alerts := make([]AlertResponse, len(events))
	for i, e := range events {
		alerts[i] = mapEventToAlertResponse(e)
	}

	return alerts, nil
}

func (s *service) GetAlertStats() (*AlertStatsResponse, error) {
	total, critical, err := s.repo.GetAlertStats()
	if err != nil {
		return nil, err
	}

	return &AlertStatsResponse{
		TotalAlerts:    total,
		CriticalAlerts: critical,
	}, nil
}

func mapEventToAlertResponse(e models.Event) AlertResponse {
	return AlertResponse{
		ID:        e.ID,
		ClientID:  e.ClientID,
		EventType: e.EventType,
		RiskScore: e.RiskScore,
		Timestamp: e.Timestamp,
	}
}
