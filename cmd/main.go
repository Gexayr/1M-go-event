package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"go-event-registration/configs"
	"go-event-registration/internal/event"
	"go-event-registration/pkg/alert"
	"go-event-registration/pkg/db"
	"go-event-registration/pkg/middleware"
)

func main() {
	cfg := configs.LoadConfig()

	database := db.Init(cfg)

	// Initialize alert module
	alert.Init(cfg.TelegramBotToken, cfg.TelegramChatID)

	r := gin.New()
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorHandlingMiddleware())

	r.POST("/events", event.RegisterEventHandler(database))

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
