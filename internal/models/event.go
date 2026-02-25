package models

import (
	"time"

	"gorm.io/datatypes"
)

// Event represents the database table structure
type Event struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ClientID  string         `gorm:"not null" json:"client_id"`
	EventType string         `gorm:"not null" json:"event_type"`
	Timestamp time.Time      `gorm:"not null" json:"timestamp"`
	Metadata  datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
	RiskScore int            `gorm:"default:0" json:"risk_score"`
}
