package reports

import "time"

type ReportListItem struct {
	ID          uint      `json:"id"`
	PeriodStart string    `json:"period_start"`
	PeriodEnd   string    `json:"period_end"`
	GeneratedAt time.Time `json:"generated_at"`
}

type ListReportsResponse struct {
	Data []ReportListItem `json:"data"`
}

type ReportDetailsResponse struct {
	ID          uint      `json:"id"`
	PeriodStart string    `json:"period_start"`
	PeriodEnd   string    `json:"period_end"`
	GeneratedAt time.Time `json:"generated_at"`
	Content     string    `json:"content"`
}

type GenerateReportRequest struct {
	From string `json:"from" binding:"required"`
	To   string `json:"to" binding:"required"`
}

type GenerateReportResponse struct {
	Message string `json:"message"`
}
