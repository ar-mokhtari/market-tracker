package main

import (
	"log"
	"net/http"
	"time"

	db "github.com/ar-mokhtari/market-tracker/adapter/storage"
	config "github.com/ar-mokhtari/market-tracker/config"
	delivery "github.com/ar-mokhtari/market-tracker/delivery/http"
	"github.com/rs/cors"
)

func main() {
	// 1. Initialize Configuration
	cfg := config.Init()

	// 2. Initialize Database
	database := db.Init(cfg.DBDSN)
	defer database.Close()

	// 3. Initialize Delivery
	mux := delivery.Init(database, cfg)

	// 4. CORS Setup
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap mux with CORS handler
	handler := c.Handler(mux)

	// 5. Start Server with Handler
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	currentTime := time.Now().Format("2006/01/02 15:04:05")
	log.Printf("%s 🚀 Server starting on port %s", currentTime, cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failure: %v", err)
	}
}
