package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go-event-registration/internal/models"
)

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var (
	botToken string
	chatID   string
	client   = &http.Client{
		Timeout: 10 * time.Second,
	}
)

// Init sets the configuration for the telegram alert module
func Init(token, id string) {
	botToken = token
	chatID = id
}

// SendHighRiskAlert sends a notification to Telegram for events with high risk scores
func SendHighRiskAlert(event models.Event) error {
	if botToken == "" || chatID == "" {
		return fmt.Errorf("telegram bot token or chat ID not configured")
	}

	message := fmt.Sprintf(
		"🚨 *HIGH RISK EVENT*\n\n*Client:* %s\n*Type:* %s\n*Score:* %d\n*Time:* %s",
		event.ClientID,
		event.EventType,
		event.RiskScore,
		event.Timestamp.Format(time.RFC3339),
	)

	payload := telegramPayload{
		ChatID:    chatID,
		Text:      message,
		ParseMode: "Markdown",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram payload: %w", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send telegram alert: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Description string `json:"description"`
		}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("telegram API error (status %d): %s", resp.StatusCode, errResp.Description)
	}

	log.Printf("Telegram alert sent for event ID %d (Risk Score: %d)", event.ID, event.RiskScore)
	return nil
}
