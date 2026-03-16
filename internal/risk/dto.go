package risk

import "go-event-registration/internal/models"

// ClientStats represents statistics for a single client
type ClientStats struct {
	ClientID       string  `json:"client_id"`
	TotalEvents    int     `json:"total_events"`
	HighRiskEvents int     `json:"high_risk_events"`
	AvgRisk        float64 `json:"avg_risk"`
}

// ClientListResponse is the response structure for GET /api/clients
type ClientListResponse struct {
	Data []ClientStats `json:"data"`
}

// ClientProfileResponse is the response structure for GET /api/clients/:id
type ClientProfileResponse struct {
	ClientStats
	RecentEvents []models.Event `gorm:"-" json:"recent_events"`
}
