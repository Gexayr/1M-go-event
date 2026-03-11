package dashboard

import (
	"context"
	"go-event-registration/internal/models"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type EventsResponse struct {
	Total      int64          `json:"total"`
	TotalPages int64          `json:"total_pages"`
	Data       []models.Event `json:"data"`
}

func (s *Service) GetEvents(ctx context.Context, filters EventFilters) (*EventsResponse, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	events, total, err := s.repo.GetFilteredEvents(ctx, filters)
	if err != nil {
		return nil, err
	}

	totalPages := int64(0)
	if filters.Limit > 0 {
		totalPages = (total + int64(filters.Limit) - 1) / int64(filters.Limit)
	}

	return &EventsResponse{
		Total:      total,
		TotalPages: totalPages,
		Data:       events,
	}, nil
}

func (s *Service) GetStats(ctx context.Context) (*DashboardStats, error) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.GetDashboardStats(ctx)
}

func (s *Service) GetRiskOverTime(ctx context.Context, filters RiskOverTimeFilters) (*RiskOverTimeResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := s.repo.GetRiskOverTime(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &RiskOverTimeResponse{Data: data}, nil
}

func (s *Service) GetRiskDistribution(ctx context.Context) (*RiskDistributionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := s.repo.GetRiskDistribution(ctx)
	if err != nil {
		return nil, err
	}

	return &RiskDistributionResponse{Data: data}, nil
}

func (s *Service) GetEventsPerClient(ctx context.Context) (*EventsPerClientResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	data, err := s.repo.GetEventsPerClient(ctx)
	if err != nil {
		return nil, err
	}

	return &EventsPerClientResponse{Data: data}, nil
}
