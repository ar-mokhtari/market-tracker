// Package server (db) is about setup and start database connection
package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB

// DBInit initializes and opens MySQL database connection
func DBInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	var err error
	Conn, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ cannot open database: %v", err)
	}
	if err = Conn.Ping(); err != nil {
		log.Fatalf("❌ cannot connect to database: %v", err)
	}

	log.Println("✅ database connected")

}
