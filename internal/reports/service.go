package reports

import (
	"context"
	"fmt"
	"go-event-registration/internal/report"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListReports(ctx context.Context) ([]ReportListItem, error) {
	reports, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	var items []ReportListItem
	for _, r := range reports {
		items = append(items, ReportListItem{
			ID:          r.ID,
			PeriodStart: r.PeriodStart.Format("2006-01-02"),
			PeriodEnd:   r.PeriodEnd.Format("2006-01-02"),
			GeneratedAt: r.GeneratedAt,
		})
	}
	return items, nil
}

func (s *Service) GetReport(ctx context.Context, id uint) (*ReportDetailsResponse, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ReportDetailsResponse{
		ID:          r.ID,
		PeriodStart: r.PeriodStart.Format("2006-01-02"),
		PeriodEnd:   r.PeriodEnd.Format("2006-01-02"),
		GeneratedAt: r.GeneratedAt,
		Content:     r.Content,
	}, nil
}

func (s *Service) GenerateReport(ctx context.Context, from, to string) error {
	startTime, err := time.Parse("2006-01-02", from)
	if err != nil {
		return fmt.Errorf("invalid from date: %w", err)
	}
	endTime, err := time.Parse("2006-01-02", to)
	if err != nil {
		return fmt.Errorf("invalid to date: %w", err)
	}

	// Adjust endTime to end of day
	endTime = endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db, err := s.repo.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	return report.GenerateReport(ctx, db, startTime, endTime)
}
