package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	handler "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Reverting to the simplest loading method
	_ = godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	baseURL := os.Getenv("API_BASE_URL")
	// Replace these keys with EXACTLY what you have in your .env
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST") // e.g., 127.0.0.1:3306

	// Constructing DSN manually if you don't have a single DB_DSN
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)

	// 2. Database Initialization
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, apiKey, baseURL)
	h := handler.NewPriceHandler(uc)

	go startBackgroundWorker(uc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	port := ":8080"
	fmt.Printf("2025/12/21 %s 🚀 Server starting on port %s\n", time.Now().Format("15:04:05"), port)

	server := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func startBackgroundWorker(uc *usecase.PriceUseCase) {
	_ = uc.FetchFromExternal()

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Printf("⏰ %s: Scheduled fetch started...\n", time.Now().Format("15:04:05"))
		if err := uc.FetchFromExternal(); err != nil {
			log.Printf("Worker failed: %v", err)
		} else {
			fmt.Println("Worker success")
		}
	}
}
