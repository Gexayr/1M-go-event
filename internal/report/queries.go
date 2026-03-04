package report

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AggregatedData struct {
	TotalEvents           int64
	HighRiskEvents        int64
	AvgRiskScore          float64
	TopRiskyClients       []ClientRisk
	EventsByType          map[string]int64
	PrevPeriodEvents      int64
	PrevPeriodHighRisk    int64
	EventChangePercent    float64
	HighRiskChangePercent float64
}

type ClientRisk struct {
	ClientID      string
	HighRiskCount int64
}

func CollectAggregatedData(ctx context.Context, db *sql.DB, start, end time.Time) (AggregatedData, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	data := AggregatedData{
		EventsByType: make(map[string]int64),
	}

	// 1. Total events, High risk events, Avg risk score
	query1 := `
		SELECT 
			COUNT(*), 
			COUNT(*) FILTER (WHERE risk_score >= 70), 
			COALESCE(AVG(risk_score), 0)
		FROM events 
		WHERE timestamp >= $1 AND timestamp < $2`

	err := db.QueryRowContext(ctx, query1, start, end).Scan(&data.TotalEvents, &data.HighRiskEvents, &data.AvgRiskScore)
	if err != nil {
		return data, fmt.Errorf("failed to fetch basic stats: %w", err)
	}

	// 2. Top 5 risky clients
	query2 := `
		SELECT client_id, COUNT(*) as high_risk_count
		FROM events
		WHERE risk_score >= 70 AND timestamp >= $1 AND timestamp < $2
		GROUP BY client_id
		ORDER BY high_risk_count DESC
		LIMIT 5`

	rows, err := db.QueryContext(ctx, query2, start, end)
	if err != nil {
		return data, fmt.Errorf("failed to fetch top risky clients: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cr ClientRisk
		if err := rows.Scan(&cr.ClientID, &cr.HighRiskCount); err != nil {
			return data, err
		}
		data.TopRiskyClients = append(data.TopRiskyClients, cr)
	}

	// 3. Events grouped by type
	query3 := `
		SELECT event_type, COUNT(*)
		FROM events
		WHERE timestamp >= $1 AND timestamp < $2
		GROUP BY event_type`

	rows3, err := db.QueryContext(ctx, query3, start, end)
	if err != nil {
		return data, fmt.Errorf("failed to fetch event types: %w", err)
	}
	defer rows3.Close()

	for rows3.Next() {
		var et string
		var count int64
		if err := rows3.Scan(&et, &count); err != nil {
			return data, err
		}
		data.EventsByType[et] = count
	}

	// 4. Comparison vs previous period
	duration := end.Sub(start)
	prevStart := start.Add(-duration)
	prevEnd := start

	query4 := `
		SELECT 
			COUNT(*), 
			COUNT(*) FILTER (WHERE risk_score >= 70)
		FROM events 
		WHERE timestamp >= $1 AND timestamp < $2`

	err = db.QueryRowContext(ctx, query4, prevStart, prevEnd).Scan(&data.PrevPeriodEvents, &data.PrevPeriodHighRisk)
	if err != nil {
		return data, fmt.Errorf("failed to fetch previous period stats: %w", err)
	}

	if data.PrevPeriodEvents > 0 {
		data.EventChangePercent = float64(data.TotalEvents-data.PrevPeriodEvents) / float64(data.PrevPeriodEvents) * 100
	}
	if data.PrevPeriodHighRisk > 0 {
		data.HighRiskChangePercent = float64(data.HighRiskEvents-data.PrevPeriodHighRisk) / float64(data.PrevPeriodHighRisk) * 100
	}

	return data, nil
}
