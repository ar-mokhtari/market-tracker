// main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	"github.com/ar-mokhtari/market-tracker/config/env"
	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/server"
	"github.com/ar-mokhtari/market-tracker/usecase"
)

func main() {
	// Load environment
	env.Init()

	cfg, err := env.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	server.DBInit()

	// Create tables
	if err := mysql.CreatePricesTable(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Initialize repository and use case
	repo := mysql.NewRepository(db)
	priceUseCase := usecase.NewUseCase(repo, env.AuthKey)

	// Fetch initial market data
	log.Println("Fetching initial market data...")
	ctx := context.Background()
	if err := priceUseCase.FetchAndSaveMarketData(ctx); err != nil {
		log.Printf("Warning: Failed to fetch initial market data: %v", err)
	} else {
		log.Println("Initial market data fetched successfully")
	}

	// Start periodic data fetching (every 5 minutes)
	go startPeriodicFetch(priceUseCase, 5*time.Minute)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Load API prefix
	prefix := env.LoadPrefix()
	apiGroup := e.Group(prefix.APIPrefix + prefix.APIVersion)

	// Register routes
	handler := v1.NewHandler(priceUseCase)
	handler.RegisterRoutes(apiGroup)

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status": "healthy",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	port := cfg.Service.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startPeriodicFetch(useCase *usecase.UseCase, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		log.Println("Fetching market data...")
		if err := useCase.FetchAndSaveMarketData(ctx); err != nil {
			log.Printf("Error fetching market data: %v", err)
		} else {
			log.Println("Market data updated successfully")
		}
	}
}
