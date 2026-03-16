package risk

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

func (h *Handler) GetClientsStats(c *gin.Context) {
	stats, err := h.service.GetClientsStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ClientListResponse{Data: stats})
}

func (h *Handler) GetClientProfile(c *gin.Context) {
	clientID := c.Param("id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client id is required"})
		return
	}
	profile, err := h.service.GetClientProfile(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/api/clients", h.GetClientsStats)
	r.GET("/api/clients/:id", h.GetClientProfile)
}
