package config

import (
	"log"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	HTTPPort       string
	InternalSecret string
	AreaServiceURL string
	LogAllRequests bool
}

func Load() *Config {
	cfg := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "cron_service_db"),
		HTTPPort:       getEnv("SERVER_PORT", "8086"),
		InternalSecret: getEnv("INTERNAL_SECRET", "secret"),
		AreaServiceURL: getEnv("AREA_SERVICE_URL", "http://gateway:8080/area_area_api"),
		LogAllRequests: getEnv("LOG_ALL_REQUESTS", "false") == "true",
	}

	log.Println("Configuration loaded successfully")
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
