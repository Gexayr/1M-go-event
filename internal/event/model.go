package event

import (
	"go-event-registration/internal/models"
	"time"
)

// Event represents the database table structure
type Event = models.Event

// RegisterEventRequest represents the incoming JSON body
type RegisterEventRequest struct {
	ClientID  string                 `json:"client_id" binding:"required"`
	EventType string                 `json:"event_type" binding:"required"`
	Timestamp *string                `json:"timestamp"` // ISO string
	Metadata  map[string]interface{} `json:"metadata"`
}

// RegisterEventResponse represents the success response
type RegisterEventResponse struct {
	EventID   uint      `json:"event_id"`
	Timestamp time.Time `json:"timestamp"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
