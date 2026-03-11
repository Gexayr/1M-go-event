package dashboard

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetEvents(c *gin.Context) {
	filters := EventFilters{
		ClientID: c.Query("client_id"),
		Limit:    50,
		Offset:   0,
	}

	if minScoreStr := c.Query("min_score"); minScoreStr != "" {
		if minScore, err := strconv.Atoi(minScoreStr); err == nil {
			filters.MinScore = &minScore
		}
	}

	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filters.From = &from
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filters.To = &to
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	res, err := h.service.GetEvents(c.Request.Context(), filters)
	if err != nil {
		log.Printf("[Dashboard] Error fetching events: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		log.Printf("[Dashboard] Error fetching stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch dashboard statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetRiskOverTime(c *gin.Context) {
	filters := RiskOverTimeFilters{}

	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse("2006-01-02", fromStr); err == nil {
			filters.From = &from
		} else if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filters.From = &from
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse("2006-01-02", toStr); err == nil {
			filters.To = &to
		} else if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filters.To = &to
		}
	}

	res, err := h.service.GetRiskOverTime(c.Request.Context(), filters)
	if err != nil {
		log.Printf("[Dashboard] Error fetching risk over time: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch risk over time"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetRiskDistribution(c *gin.Context) {
	res, err := h.service.GetRiskDistribution(c.Request.Context())
	if err != nil {
		log.Printf("[Dashboard] Error fetching risk distribution: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch risk distribution"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetEventsPerClient(c *gin.Context) {
	res, err := h.service.GetEventsPerClient(c.Request.Context())
	if err != nil {
		log.Printf("[Dashboard] Error fetching events per client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events per client"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	api := r.Group("/api")
	{
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/events", h.GetEvents)
			dashboard.GET("/stats", h.GetStats)
			dashboard.GET("/risk-over-time", h.GetRiskOverTime)
			dashboard.GET("/risk-distribution", h.GetRiskDistribution)
			dashboard.GET("/events-per-client", h.GetEventsPerClient)
		}
	}
}
