package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	TelegramToken    string
	TelegramChatID   string
	ServerPort       string
	ShutdownTimeout  time.Duration
	ScrapingInterval time.Duration
	PriceHistoryDays int
	WorkerPoolSize   int
}

func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	shutdownTimeout, _ := strconv.Atoi(getEnv("SHUTDOWN_TIMEOUT", "30"))
	scrapingInterval, _ := strconv.Atoi(getEnv("SCRAPING_INTERVAL", "3600")) // 1 hour default
	priceHistoryDays, _ := strconv.Atoi(getEnv("PRICE_HISTORY_DAYS", "30"))

	workerPoolSize, _ := strconv.Atoi(getEnv("WORKER_POOL_SIZE", "50"))

	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/price_watcher?sslmode=disable"),
		TelegramToken:    getEnv("TELEGRAM_TOKEN", ""),
		TelegramChatID:   getEnv("TELEGRAM_CHAT_ID", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		ShutdownTimeout:  time.Duration(shutdownTimeout) * time.Second,
		ScrapingInterval: time.Duration(scrapingInterval) * time.Second,
		PriceHistoryDays: priceHistoryDays,
		WorkerPoolSize:   workerPoolSize,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
