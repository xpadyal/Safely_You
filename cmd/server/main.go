// Package main is the entry point for the Safely You device monitoring server.
// It initializes the application, loads device data, and starts the HTTP server.
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/xpadyal/Safely_You/internal/config"
	"github.com/xpadyal/Safely_You/internal/handlers"
	"github.com/xpadyal/Safely_You/internal/loader"
	"github.com/xpadyal/Safely_You/internal/store"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize store
	storeInstance := store.NewStore()

	// Load devices from CSV file
	if err := loader.LoadDevicesFromCSV("devices.csv", storeInstance); err != nil {
		log.Fatalf("Failed to load devices: %v", err)
	}

	// Set up the web server
	r := gin.Default()

	// Health check
	r.GET("/health", handlers.HealthHandler)

	// API routes for devices
	api := r.Group("/api/v1")
	devices := api.Group("/devices")
	{
		devices.POST("/:device_id/heartbeat", handlers.PostHeartbeatHandler(storeInstance))
		devices.POST("/:device_id/stats", handlers.PostStatsHandler(storeInstance))
		devices.GET("/:device_id/stats", handlers.GetStatsHandler(storeInstance))
	}

	log.Printf("listening on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}
}
