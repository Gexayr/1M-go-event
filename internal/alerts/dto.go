package alerts

import "time"

type AlertResponse struct {
	ID        uint      `json:"id"`
	ClientID  string    `json:"client_id"`
	EventType string    `json:"event_type"`
	RiskScore int       `json:"risk_score"`
	Timestamp time.Time `json:"timestamp"`
}

type AlertStatsResponse struct {
	TotalAlerts    int64 `json:"total_alerts"`
	CriticalAlerts int64 `json:"critical_alerts"`
}

type AlertsListResponse struct {
	Data []AlertResponse `json:"data"`
}
