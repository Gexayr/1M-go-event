package configs

import "os"

type Config struct {
	DatabaseURL      string
	Port             string
	TelegramBotToken string
	TelegramChatID   string
}

func LoadConfig() *Config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=events port=5432 sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL:      dsn,
		Port:             port,
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	}
}
