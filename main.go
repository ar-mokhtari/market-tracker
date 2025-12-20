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
	godotenv.Load()

	// استفاده از متغیرهای محیطی با دقت بالا
	dbUser := os.Getenv("DB_USER") // market_user
	dbPass := os.Getenv("DB_PASS") // market_secure_password_2024
	dbHost := os.Getenv("DB_HOST") // localhost یا 127.0.0.1
	dbPort := os.Getenv("DB_PORT") // 3306
	dbName := os.Getenv("DB_NAME") // market_tracker

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("خطا در باز کردن دیتابیس: ", err)
	}

	// تست واقعی اتصال قبل از شروع
	err = db.Ping()
	if err != nil {
		log.Fatal("دیتابیس در دسترس نیست: ", err)
	}

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
);`

	_, err = db.Exec(migration)
	if err != nil {
		log.Fatal("خطا در ساخت جداول: ", err)
	}

	apiKey := os.Getenv("API_KEY") // خواندن از .env
	repo := mysql.NewRepository(db)
	uc := usecase.NewPriceUseCase(repo, apiKey)

	h := v1.NewHandler(uc)

	// شروع دریافت خودکار
	go func() {
		for {
			uc.FetchFromExternal()
			time.Sleep(5 * time.Minute)
		}
	}()

	// مسیرها (Routes)
	http.HandleFunc("/api/v1/prices", h.ListPrices)
	http.HandleFunc("/api/v1/prices/fetch", h.ManualFetch)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
