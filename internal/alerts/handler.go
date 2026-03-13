package alerts

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAlerts(c *gin.Context) {
	alerts, err := h.service.GetAlerts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AlertsListResponse{
		Data: alerts,
	})
}

func (h *Handler) GetAlertStats(c *gin.Context) {
	stats, err := h.service.GetAlertStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		api.GET("/alerts", h.GetAlerts)
		api.GET("/alerts/stats", h.GetAlertStats)
	}
}
