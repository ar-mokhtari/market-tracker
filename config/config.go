// Package config manages application configuration from environment variables.
package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	dbUser        string
	dbPass        string
	dbHost        string
	dbName        string
	DBDSN         string
	Port          string
	APIKey        string
	BaseURL       string
	FetchInterval int
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := os.Getenv(key)
	if value := valueStr; value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultVal
}

func Init() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	cfg := &Config{
		Port:    os.Getenv("PORT"),
		APIKey:  os.Getenv("API_KEY"),
		BaseURL: os.Getenv("API_BASE_URL"),
		dbUser:  os.Getenv("DB_USER"),
		dbPass:  os.Getenv("DB_PASS"),
		dbName:  os.Getenv("DB_NAME"),
		dbHost:  os.Getenv("DB_HOST"),
	}

	cfg.DBDSN = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.dbUser, cfg.dbPass, cfg.dbHost, cfg.dbName)

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DBDSN == "" {
		log.Fatal("Critical: DB_DSN is not set in environment")
	}

	cfg.FetchInterval = getEnvAsInt("FETCH_INTERVAL_MINUTES", 10)

	return cfg
}
