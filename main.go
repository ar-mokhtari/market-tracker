package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ar-mokhtari/market-tracker/adapter/storage/mysql"
	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/usecase"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Information: .env file not found, using system env: %v", err)
	}

	db := initDB()
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	runMigrations(db)

	apiKey := os.Getenv("API_KEY")
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, apiKey)
	handler := v1.NewHandler(uc)

	go startWorker(uc)

	http.HandleFunc("/api/v1/prices", handler.ListPrices)
	http.HandleFunc("/api/v1/prices/fetch", handler.ManualFetch)
	http.HandleFunc("/api/v1/prices/timeline", handler.GetTimeline)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Captain Aliz, Server is running on port %s...\n", port)

	server := &http.Server{
		Addr:              ":" + port,
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func initDB() *sql.DB {
	// Added multiStatements=true to fix the Migration error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("DB Open Error: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("DB Connection Error: %v", err)
	}
	return db
}

func runMigrations(db *sql.DB) {
	migration := `
	CREATE TABLE IF NOT EXISTS prices (
		id INT AUTO_INCREMENT PRIMARY KEY,
		date VARCHAR(20),
		time VARCHAR(20),
		symbol VARCHAR(50) NOT NULL,
		name_fa VARCHAR(100),
		price VARCHAR(50),
		unit VARCHAR(20),
		type VARCHAR(20) NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		UNIQUE KEY unique_symbol_type (symbol, type)
	);
	CREATE TABLE IF NOT EXISTS price_history (
		id INT AUTO_INCREMENT PRIMARY KEY,
		symbol VARCHAR(50),
		price VARCHAR(50),
		type VARCHAR(20),
		recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(migration); err != nil {
		log.Fatalf("Migration Error: %v", err)
	}
}

func startWorker(uc *usecase.PriceUseCase) {
	if err := uc.FetchFromExternal(); err != nil {
		log.Printf("Initial fetch error: %v", err)
	}

	ticker := time.NewTicker(60 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := uc.FetchFromExternal(); err != nil {
			log.Printf("Worker Error: %v", err)
		}
	}
}
