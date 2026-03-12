package events

import (
	"encoding/json"
	"time"
)

type EventDetailsResponse struct {
	ID        uint            `json:"id"`
	ClientID  uint            `json:"client_id"`
	EventType string          `json:"event_type"`
	RiskScore int             `json:"risk_score"`
	Metadata  json.RawMessage `json:"metadata"`
	Timestamp time.Time       `json:"timestamp"`
}
