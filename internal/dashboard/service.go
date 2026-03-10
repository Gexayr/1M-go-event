package dashboard

import (
	"context"
	"go-event-registration/internal/models"
	"time"
)

type Service struct {
	queries *Queries
}

func NewService(queries *Queries) *Service {
	return &Service{queries: queries}
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

	events, total, err := s.queries.GetFilteredEvents(ctx, filters)
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

	return s.queries.GetDashboardStats(ctx)
}
