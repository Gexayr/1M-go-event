package risk

import (
	"encoding/json"
	"fmt"
	"go-event-registration/internal/models"
	"go-event-registration/internal/rules"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// CalculateRisk evaluates the risk score for a given event.
func CalculateRisk(evt models.Event, db *gorm.DB) (int, error) {
	score := 0

	// Load enabled rules from database
	var activeRules []rules.RiskRule
	if err := db.Where("enabled = ?", true).Find(&activeRules).Error; err != nil {
		return 0, err
	}

	var metadata map[string]interface{}
	if len(evt.Metadata) > 0 {
		if err := json.Unmarshal(evt.Metadata, &metadata); err != nil {
			// Log error but continue with other rules
			fmt.Printf("failed to unmarshal metadata: %v\n", err)
		}
	}

	// Evaluate dynamic rules
	for _, rule := range activeRules {
		if evt.EventType == rule.EventType {
			if metadata == nil {
				// If no metadata but rule expects a field, skip or fail?
				// Usually skip as condition cannot be met.
				continue
			}

			val, ok := metadata[rule.ConditionField]
			if !ok {
				continue
			}

			if evaluateCondition(val, rule.ConditionOperator, rule.ConditionValue) {
				score += rule.Score
			}
		}
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

func evaluateCondition(actual interface{}, operator string, target string) bool {
	// For simplicity, we'll try to convert everything to float64 if it's a number
	// or compare as strings.

	actualStr := fmt.Sprintf("%v", actual)

	switch operator {
	case ">":
		a, err1 := strconv.ParseFloat(actualStr, 64)
		t, err2 := strconv.ParseFloat(target, 64)
		if err1 == nil && err2 == nil {
			return a > t
		}
	case "<":
		a, err1 := strconv.ParseFloat(actualStr, 64)
		t, err2 := strconv.ParseFloat(target, 64)
		if err1 == nil && err2 == nil {
			return a < t
		}
	case ">=":
		a, err1 := strconv.ParseFloat(actualStr, 64)
		t, err2 := strconv.ParseFloat(target, 64)
		if err1 == nil && err2 == nil {
			return a >= t
		}
	case "<=":
		a, err1 := strconv.ParseFloat(actualStr, 64)
		t, err2 := strconv.ParseFloat(target, 64)
		if err1 == nil && err2 == nil {
			return a <= t
		}
	case "==":
		return actualStr == target
	case "!=":
		return actualStr != target
	}

	return false
}
