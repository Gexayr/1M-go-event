package risk

import (
	"encoding/json"
	"go-event-registration/internal/models"
	"time"

	"gorm.io/gorm"
)

// CalculateRisk evaluates the risk score for a given event.
func CalculateRisk(evt models.Event, db *gorm.DB) (int, error) {
	score := 0

	// Rule 1: login_failed event type
	if evt.EventType == "login_failed" {
		score += 30
	}

	// Rule 2: more than 5 login_failed events for same client in last 10 minutes
	if evt.EventType == "login_failed" {
		var count int64
		tenMinutesAgo := time.Now().Add(-10 * time.Minute)
		err := db.Model(&models.Event{}).
			Where("client_id = ? AND event_type = ? AND timestamp > ?", evt.ClientID, "login_failed", tenMinutesAgo).
			Count(&count).Error
		if err != nil {
			return 0, err
		}

		if count > 5 {
			score += 50
		}
	}

	// Rule 3: withdrawal AND amount > 10000
	if evt.EventType == "withdrawal" {
		var metadata map[string]interface{}
		if err := json.Unmarshal(evt.Metadata, &metadata); err == nil {
			if amount, ok := metadata["amount"].(float64); ok && amount > 10000 {
				score += 70
			}
		}
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	return score, nil
}
