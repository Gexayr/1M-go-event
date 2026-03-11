package dashboard

import (
	"context"
	"go-event-registration/internal/models"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

type EventFilters struct {
	ClientID string
	MinScore *int
	From     *time.Time
	To       *time.Time
	Limit    int
	Offset   int
}

func (r *Repository) GetFilteredEvents(ctx context.Context, filters EventFilters) ([]models.Event, int64, error) {
	var events []models.Event
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Event{})

	if filters.ClientID != "" {
		query = query.Where("client_id = ?", filters.ClientID)
	}
	if filters.MinScore != nil {
		query = query.Where("risk_score >= ?", *filters.MinScore)
	}
	if filters.From != nil {
		query = query.Where("timestamp >= ?", *filters.From)
	}
	if filters.To != nil {
		query = query.Where("timestamp <= ?", *filters.To)
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch data with pagination
	if err := query.Limit(filters.Limit).Offset(filters.Offset).Order("timestamp DESC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

type DashboardStats struct {
	TotalEvents     int64               `json:"total_events"`
	HighRiskEvents  int64               `json:"high_risk_events"`
	AvgRiskScore    float64             `json:"avg_risk_score"`
	TopRiskyClients []TopRiskyClientRow `json:"top_risky_clients"`
}

type TopRiskyClientRow struct {
	ClientID      string `json:"client_id"`
	HighRiskCount int64  `json:"high_risk_count"`
}

func (r *Repository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	var stats DashboardStats

	// Total events
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Count(&stats.TotalEvents).Error; err != nil {
		return nil, err
	}

	// High risk events (assuming score > 70 is high risk, or whatever defines "high risk")
	// The requirement doesn't specify the threshold, common threshold is 70 or 80.
	// I'll use 70 for now.
	highRiskThreshold := 70
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Where("risk_score > ?", highRiskThreshold).Count(&stats.HighRiskEvents).Error; err != nil {
		return nil, err
	}

	// Avg risk score
	if err := r.db.WithContext(ctx).Model(&models.Event{}).Select("COALESCE(AVG(risk_score), 0)").Scan(&stats.AvgRiskScore).Error; err != nil {
		return nil, err
	}

	// Top risky clients
	if err := r.db.WithContext(ctx).Model(&models.Event{}).
		Select("client_id, COUNT(*) as high_risk_count").
		Where("risk_score > ?", highRiskThreshold).
		Group("client_id").
		Order("high_risk_count DESC").
		Limit(5).
		Scan(&stats.TopRiskyClients).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *Repository) GetRiskOverTime(ctx context.Context, filters RiskOverTimeFilters) ([]RiskOverTimeData, error) {
	var results []RiskOverTimeData

	query := r.db.WithContext(ctx).Model(&models.Event{}).
		Select("DATE(timestamp) as date, AVG(risk_score) as avg_risk")

	if filters.From != nil {
		query = query.Where("timestamp >= ?", *filters.From)
	}
	if filters.To != nil {
		query = query.Where("timestamp <= ?", *filters.To)
	}

	err := query.Group("DATE(timestamp)").
		Order("date").
		Scan(&results).Error

	return results, err
}

func (r *Repository) GetRiskDistribution(ctx context.Context) ([]RiskDistributionData, error) {
	var results []RiskDistributionData

	// Use a CASE statement to group risk scores into ranges
	sql := `
		SELECT 
			ranges.range_name as range,
			COUNT(events.id) as count
		FROM (
			SELECT '0-20' as range_name, 0 as min_score, 20 as max_score
			UNION ALL SELECT '21-40', 21, 40
			UNION ALL SELECT '41-60', 41, 60
			UNION ALL SELECT '61-80', 61, 80
			UNION ALL SELECT '81-100', 81, 100
		) as ranges
		LEFT JOIN events ON events.risk_score BETWEEN ranges.min_score AND ranges.max_score
		GROUP BY ranges.range_name, ranges.min_score
		ORDER BY ranges.min_score
	`

	err := r.db.WithContext(ctx).Raw(sql).Scan(&results).Error
	return results, err
}

func (r *Repository) GetEventsPerClient(ctx context.Context) ([]EventsPerClientData, error) {
	var results []EventsPerClientData

	err := r.db.WithContext(ctx).Model(&models.Event{}).
		Select("client_id, COUNT(*) as events").
		Group("client_id").
		Order("events DESC").
		Limit(10).
		Scan(&results).Error

	return results, err
}
