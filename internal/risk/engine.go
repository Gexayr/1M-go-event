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
	if evt.EventType == "failed_login" {
		score += 30
	}

	// Rule 2: more than 5 failed_login events for same client in last 10 minutes
	if evt.EventType == "failed_login" {
		var count int64
		tenMinutesAgo := time.Now().Add(-10 * time.Minute)
		err := db.Model(&models.Event{}).
			Where("client_id = ? AND event_type = ? AND timestamp > ?", evt.ClientID, "failed_login", tenMinutesAgo).
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

	// Rule 4: deposit AND amount > 5000
	if evt.EventType == "deposit" {
		var metadata map[string]interface{}
		if err := json.Unmarshal(evt.Metadata, &metadata); err == nil {
			if amount, ok := metadata["amount"].(float64); ok && amount > 5000 {
				score += 40
			}
		}
	}

	// Rule 5: more than 20 events for same client in last 5 minutes
	var eventCount int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	err := db.Model(&models.Event{}).
		Where("client_id = ? AND timestamp > ?", evt.ClientID, fiveMinutesAgo).
		Count(&eventCount).Error

	if err == nil && eventCount > 20 {
		score += 35
	}

	// Rule 6: more than 3 withdrawals in 10 minutes
	if evt.EventType == "withdrawal" {

		var withdrawalCount int64
		tenMinutesAgo := time.Now().Add(-10 * time.Minute)

		err := db.Model(&models.Event{}).
			Where("client_id = ? AND event_type = ? AND timestamp > ?", evt.ClientID, "withdrawal", tenMinutesAgo).
			Count(&withdrawalCount).Error

		if err == nil && withdrawalCount > 3 {
			score += 60
		}
	}

	// Rule 7: activity between 2AM–5AM
	hour := time.Now().Hour()

	if hour >= 2 && hour <= 5 {
		score += 15
	}

	// Rule 8: client had high risk events recently
	var riskyEvents int64
	oneHourAgo := time.Now().Add(-1 * time.Hour)

	err = db.Model(&models.Event{}).
		Where("client_id = ? AND risk_score > ? AND timestamp > ?", evt.ClientID, 70, oneHourAgo).
		Count(&riskyEvents).Error

	if err == nil && riskyEvents > 3 {
		score += 45
	}

	// Cap score at 100
	if score > 100 {
		score = 100
	}

	return score, nil
}
