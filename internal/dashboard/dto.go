package dashboard

import "time"

type RiskOverTimeData struct {
	Date    string  `json:"date"`
	AvgRisk float64 `json:"avg_risk"`
}

type RiskOverTimeResponse struct {
	Data []RiskOverTimeData `json:"data"`
}

type RiskDistributionData struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

type RiskDistributionResponse struct {
	Data []RiskDistributionData `json:"data"`
}

type EventsPerClientData struct {
	ClientID string `json:"client_id"`
	Events   int64  `json:"events"`
}

type EventsPerClientResponse struct {
	Data []EventsPerClientData `json:"data"`
}

type RiskOverTimeFilters struct {
	From *time.Time
	To   *time.Time
}
