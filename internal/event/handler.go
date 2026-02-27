package event

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"go-event-registration/pkg/alert"
	"go-event-registration/pkg/risk"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RegisterEventHandler handles the POST /events registration logic
func RegisterEventHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterEventRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}

		// Handle timestamp validation and default
		var eventTime time.Time
		if req.Timestamp == nil || *req.Timestamp == "" {
			eventTime = time.Now().UTC()
		} else {
			var err error
			eventTime, err = time.Parse(time.RFC3339, *req.Timestamp)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid timestamp format, expected ISO 8601 (RFC3339)"})
				return
			}
		}

		// Handle Metadata conversion to JSONB
		metadataJSON := datatypes.JSON("{}")
		if req.Metadata != nil {
			mBytes, err := json.Marshal(req.Metadata)
			if err != nil {
				c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to process metadata"})
				return
			}
			metadataJSON = datatypes.JSON(mBytes)
		}

		evt := Event{
			ClientID:  req.ClientID,
			EventType: req.EventType,
			Timestamp: eventTime,
			Metadata:  metadataJSON,
		}

		// Save to database
		if err := db.Create(&evt).Error; err != nil {
			log.Printf("Error saving event to database: %v", err)
			c.Error(err)
			return
		}

		// Calculate risk
		riskScore, err := risk.CalculateRisk(evt, db)
		if err != nil {
			log.Printf("Error calculating risk: %v", err)
			// We continue even if risk calculation fails, or we could return error
		} else {
			evt.RiskScore = riskScore
			if err := db.Model(&evt).Update("risk_score", riskScore).Error; err != nil {
				log.Printf("Error updating risk score: %v", err)
			}

			if riskScore >= 70 {
				log.Printf("HIGH RISK EVENT: ClientID=%s, EventType=%s, Score=%d", evt.ClientID, evt.EventType, riskScore)
				if err := alert.SendHighRiskAlert(evt); err != nil {
					log.Printf("Error sending telegram alert: %v", err)
				}
			}
		}

		c.JSON(http.StatusCreated, RegisterEventResponse{
			EventID:   evt.ID,
			Timestamp: evt.Timestamp,
		})
	}
}
