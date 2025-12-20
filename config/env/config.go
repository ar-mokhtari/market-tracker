// Package env is about config env like key(s)
package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var AuthKey string

func Init() {
	envPath := os.Getenv("ENV_PATH")
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("failed to load .env file from '%s': %v\n", envPath, err)
		}
	} else {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current working directory: %v", err)
		}
		envPath = filepath.Join(wd, ".env")
		if err := godotenv.Load(envPath); err != nil {
			log.Printf("warning: failed to load .env file from '%s': %v\n", envPath, err)
		}
	}

	AuthKey = GetConf("AUTH_KEY", "")
	if len(AuthKey) < 20 {
		panic("security error: AuthKey too short")
	}
}

func GetConf(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type Prefix struct {
	APIPrefix  string
	APIVersion string
}

func LoadPrefix() *Prefix {
	return &Prefix{
		APIPrefix:  "/api",
		APIVersion: "/v1",
	}
}
