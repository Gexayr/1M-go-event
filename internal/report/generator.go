package report

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type LLMRequest struct {
	Model    string       `json:"model"`
	Messages []LLMMessage `json:"messages"`
}

type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func GenerateReport(ctx context.Context, db *sql.DB, periodStart, periodEnd time.Time) error {
	// 1. Collect aggregated data
	data, err := CollectAggregatedData(ctx, db, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("collecting data: %w", err)
	}

	// 2. Build LLM prompt
	prompt := BuildReportPrompt(data)

	// 3. Call LLM API
	content, err := callLLMAPI(ctx, prompt)
	if err != nil {
		return fmt.Errorf("calling LLM API: %w", err)
	}

	// 4. Store result in reports table
	query := `
		INSERT INTO reports (period_start, period_end, generated_at, content)
		VALUES ($1, $2, $3, $4)`

	_, err = db.ExecContext(ctx, query, periodStart, periodEnd, time.Now(), content)
	if err != nil {
		return fmt.Errorf("storing report: %w", err)
	}

	return nil
}

func callLLMAPI(ctx context.Context, prompt string) (string, error) {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("LLM_API_KEY environment variable not set")
	}

	reqBody := LLMRequest{
		Model: "gpt-4",
		Messages: []LLMMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Using a standard endpoint, assuming OpenAI-compatible API
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API returned status: %s", resp.Status)
	}

	var llmResp LLMResponse
	if err := json.NewDecoder(resp.Body).Decode(&llmResp); err != nil {
		return "", err
	}

	if len(llmResp.Choices) == 0 {
		return "", fmt.Errorf("LLM API returned no choices")
	}

	return llmResp.Choices[0].Message.Content, nil
}

func GenerateWeeklyReport(db *sql.DB) error {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	start := end.AddDate(0, 0, -7)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	return GenerateReport(ctx, db, start, end)
}
