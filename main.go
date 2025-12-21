// Package main provides the entry point for the market-tracker service,
// managing the initialization of config, database, and delivery layers.
package main

import (
	"log"
	"net/http"
	"time"

	db "github.com/ar-mokhtari/market-tracker/adapter/storage"
	config "github.com/ar-mokhtari/market-tracker/config"
	delivery "github.com/ar-mokhtari/market-tracker/delivery/http"
)

func main() {
	// 1. Initialize Configuration
	cfg := config.Init()

	// 2. Initialize Database
	database := db.Init(cfg.DBDSN)
	defer database.Close()

	// 3. Initialize Delivery (Handlers, Usecases, and Background Workers)
	mux := delivery.Init(database, cfg)

	// 4. Start Server with Dynamic Port
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	currentTime := time.Now().Format("2006/01/02 15:04:05")
	log.Printf("%s 🚀 Server starting on port %s", currentTime, cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failure: %v", err)
	}
}
