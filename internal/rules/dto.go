package rules

import (
	"time"

	"gorm.io/gorm"
)

// RiskRule represents the database table structure
type RiskRule struct {
	ID                uint           `json:"id"`
	Name              string         `json:"rule_name"`
	Description       string         `json:"description"`
	EventType         string         `json:"event_type"`
	ConditionField    string         `json:"condition_field"`
	ConditionOperator string         `json:"condition_operator"`
	ConditionValue    string         `json:"condition_value"`
	Score             int            `json:"score"`
	Enabled           bool           `json:"enabled"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-"`
}

// CreateRuleRequest represents the POST /api/rules request body
type CreateRuleRequest struct {
	Name              string `json:"rule_name" binding:"required"`
	Description       string `json:"description"`
	EventType         string `json:"event_type" binding:"required"`
	ConditionField    string `json:"condition_field" binding:"required"`
	ConditionOperator string `json:"condition_operator" binding:"required"`
	ConditionValue    string `json:"condition_value" binding:"required"`
	Score             int    `json:"score" binding:"required"`
	Enabled           bool   `json:"enabled"`
}

// UpdateRuleRequest represents the PUT /api/rules/:id request body
type UpdateRuleRequest struct {
	Name              *string `json:"rule_name"`
	Description       *string `json:"description"`
	EventType         *string `json:"event_type"`
	ConditionField    *string `json:"condition_field"`
	ConditionOperator *string `json:"condition_operator"`
	ConditionValue    *string `json:"condition_value"`
	Score             *int    `json:"score"`
	Enabled           *bool   `json:"enabled"`
}

// RuleResponse represents a single rule in JSON responses
type RuleResponse struct {
	ID                uint      `json:"id"`
	Name              string    `json:"rule_name"`
	Description       string    `json:"description"`
	EventType         string    `json:"event_type"`
	ConditionField    string    `json:"condition_field"`
	ConditionOperator string    `json:"condition_operator"`
	ConditionValue    string    `json:"condition_value"`
	Score             int       `json:"score"`
	Enabled           bool      `json:"enabled"`
	CreatedAt         time.Time `json:"created_at"`
}

// NewRuleResponse maps RiskRule to RuleResponse
func NewRuleResponse(rule RiskRule) RuleResponse {
	return RuleResponse{
		ID:                rule.ID,
		Name:              rule.Name,
		Description:       rule.Description,
		EventType:         rule.EventType,
		ConditionField:    rule.ConditionField,
		ConditionOperator: rule.ConditionOperator,
		ConditionValue:    rule.ConditionValue,
		Score:             rule.Score,
		Enabled:           rule.Enabled,
		CreatedAt:         rule.CreatedAt,
	}
}
