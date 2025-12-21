// Package server provides initialization for database, routes, and handlers.
package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	v1 "github.com/ar-mokhtari/market-tracker/delivery/http/v1"
	"github.com/ar-mokhtari/market-tracker/usecase"
	_ "github.com/go-sql-driver/mysql" // Register MySQL driver
)

// DBInit initializes the MySQL database connection.
func DBInit() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Critical: Could not open database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Critical: Database is unreachable: %v", err)
	}

	return db
}

// HandlerInit creates a new HTTP handler instance.
func HandlerInit(uc *usecase.PriceUseCase) *v1.Handler {
	return v1.NewHandler(uc)
}
