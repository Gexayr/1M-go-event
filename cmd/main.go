package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-event-registration/configs"
	"go-event-registration/internal/alert"
	"go-event-registration/internal/alerts"
	"go-event-registration/internal/dashboard"
	"go-event-registration/internal/event"
	"go-event-registration/internal/events"
	"go-event-registration/pkg/db"
	"go-event-registration/pkg/middleware"
)

func main() {
	cfg := configs.LoadConfig()

	database := db.Init(cfg)

	// Initialize alert module
	alert.Init(cfg.TelegramBotToken, cfg.TelegramChatID)

	// Initialize Dashboard
	dashboardRepo := dashboard.NewRepository(database)
	dashboardService := dashboard.NewService(dashboardRepo)
	dashboardHandler := dashboard.NewHandler(dashboardService)

	// Initialize Events
	eventsRepo := events.NewRepository(database)
	eventsService := events.NewService(eventsRepo)
	eventsHandler := events.NewHandler(eventsService)

	// Initialize Alerts
	alertsRepo := alerts.NewRepository(database)
	alertsService := alerts.NewService(alertsRepo)
	alertsHandler := alerts.NewHandler(alertsService)

	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandlingMiddleware())

	// Enable CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.POST("/events", event.RegisterEventHandler(database))
	r.GET("/api/events/:id", eventsHandler.GetEventDetails)

	// Register Dashboard routes
	dashboard.RegisterRoutes(r, dashboardHandler)

	// Register Alerts routes
	alerts.RegisterRoutes(r, alertsHandler)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
