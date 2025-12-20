package env

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Cache    CacheConfig
	JWT      JWTConfig
	Service  ServiceConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	RootPass string
}

type CacheConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	Secret      string
	ExpirationH int
}

type ServiceConfig struct {
	Name string
	Port string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables directly")
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", ""),
			RootPass: getEnv("DB_ROOT_PASSWORD", ""),
		},
		Cache: CacheConfig{
			Host: getEnv("CACHE_HOST", ""),
			Port: getEnv("CACHE_PORT", ""),
		},
		JWT: JWTConfig{
			Secret:      getEnv("AUTH_KEY", ""),
			ExpirationH: mustAtoi(getEnv("AUTH_EXPIRATION_HOURS", "")),
		},
		Service: ServiceConfig{
			Name: getEnv("SERVICE_NAME", ""),
			Port: getEnv("SERVICE_PORT", ""),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
