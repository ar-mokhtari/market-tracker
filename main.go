// Package main is the entry point of the market-tracker application.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	"github.com/ar-mokhtari/market-tracker/delivery/http/handler"
	"github.com/ar-mokhtari/market-tracker/usecase"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	apiKey := os.Getenv("API_KEY")
	baseURL := "https://brsapi.ir/Api/Market/Gold_Currency.php"
	dbDSN := os.Getenv("DB_DSN")

	// 2. Initialize Infrastructure (MySQL)
	repo, err := mysql.NewRepository(dbDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Initialize UseCase
	uc := usecase.NewPriceUseCase(repo, apiKey, baseURL)

	// 4. Start Background Worker (Automated Fetching every 10 minutes)
	go startBackgroundWorker(uc)

	// 5. Initialize Handlers and Routes
	h := handler.NewPriceHandler(uc)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/prices/fetch", h.FetchPrices)
	mux.HandleFunc("/api/v1/prices", h.GetPrices)
	mux.HandleFunc("/api/v1/prices/timeline", h.GetTimeline)

	// 6. Start HTTP Server
	port := ":8080"
	fmt.Printf("2025/12/21 %s 🚀 Server starting on port %s\n", time.Now().Format("15:04:05"), port)

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// startBackgroundWorker runs in a goroutine to fetch data periodically.
func startBackgroundWorker(uc *usecase.PriceUseCase) {
	// Execute immediately on startup
	if err := uc.FetchFromExternal(); err != nil {
		log.Printf("Initial automation fetch failed: %v", err)
	}

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("⏰ %s: Scheduled fetch started...\n", time.Now().Format("15:04:05"))
		if err := uc.FetchFromExternal(); err != nil {
			log.Printf("❌ Scheduled fetch failed: %v", err)
		} else {
			fmt.Println("✅ Scheduled fetch successful")
		}
	}
}
