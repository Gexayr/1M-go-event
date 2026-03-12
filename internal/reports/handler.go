package reports

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListReports(c *gin.Context) {
	reports, err := h.service.ListReports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list reports"})
		return
	}

	c.JSON(http.StatusOK, ListReportsResponse{Data: reports})
}

func (h *Handler) GetReportDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid report id"})
		return
	}

	report, err := h.service.GetReport(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) GenerateReport(c *gin.Context) {
	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.GenerateReport(c.Request.Context(), req.From, req.To)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate report: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenerateReportResponse{Message: "Report generation started/completed"})
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		reports := api.Group("/reports")
		{
			reports.GET("", h.ListReports)
			reports.GET("/:id", h.GetReportDetails)
			reports.POST("/generate", h.GenerateReport)
		}
	}
}
