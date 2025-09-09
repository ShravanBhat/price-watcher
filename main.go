package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"price-watcher/config"
	"price-watcher/database"
	"price-watcher/scheduler"
	"price-watcher/server"
	"price-watcher/telegram"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Telegram bot
	tgBot, err := telegram.NewBot(cfg.TelegramToken, cfg.TelegramChatID)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	// Initialize scheduler
	sched := scheduler.NewScheduler(db, tgBot, cfg)

	// Start scheduler
	sched.Start()

	// Initialize and start HTTP server
	srv := server.NewServer(db, cfg)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	sched.Stop()
	log.Println("Server exited")
}
