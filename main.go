// Package main is the entry point of the market-tracker service.
package main

import (
	"log"
	"os"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	"github.com/ar-mokhtari/market-tracker/internal/server"
	"github.com/ar-mokhtari/market-tracker/usecase"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Info: .env file not found")
	}

	db := server.DBInit() // Assume you moved DB logic here
	defer db.Close()

	apiKey := os.Getenv("API_KEY")
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, apiKey)

	// Start the automated worker in background
	go uc.StartAutomation()

	handler := server.HandlerInit(uc)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.RouterInit(handler, port)
	log.Printf("🚀 Server starting on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
