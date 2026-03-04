package report

import (
	"encoding/json"
	"fmt"
)

func BuildReportPrompt(data AggregatedData) string {
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	prompt := fmt.Sprintf(`
You are a senior risk analyst. Generate a comprehensive risk monitoring report based on the following aggregated data:

%s

Requirements for the report:
1. Format: Markdown
2. Sections:
   - Executive Summary: Brief overview of the risk landscape.
   - Risk Overview: Analysis of total events, high-risk events, and average risk score.
   - Trends & Anomalies: Comparison with the previous period and detection of any unusual patterns.
   - Top Risky Clients: Analysis of the most problematic clients.
   - Recommendations: Actionable steps to mitigate identified risks.
3. Tone: Professional and executive-level.
4. Insights: Beyond just stating numbers, provide analysis on what these numbers mean (e.g., "High-risk events increased by X%%, suggesting a potential new threat vector").
`, string(jsonData))

	return prompt
}
