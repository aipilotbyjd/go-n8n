package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jaydeep/go-n8n/configs"
	"github.com/jaydeep/go-n8n/internal/interfaces/http/rest/v1"
	"github.com/jaydeep/go-n8n/pkg/database"
	"github.com/jaydeep/go-n8n/pkg/logger"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting n8n Clone API Server",
		"version", Version,
		"build_time", BuildTime,
	)

	// Load configuration
	cfg, err := configs.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	// Initialize Redis
	// TODO: Add Redis initialization

	// Initialize router
	router := v1.NewRouter(cfg, db, log)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		log.Info("API Server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown", "error", err)
	}

	log.Info("Server exited")
}
