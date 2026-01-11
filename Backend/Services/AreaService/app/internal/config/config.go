package config

import "os"

type Config struct {
	HTTPPort          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	AuthServiceURL    string
	ServiceServiceURL string
	AreaServiceURL    string
	InternalSecret    string
}

func Load() Config {
	return Config{
		HTTPPort:          getEnv("SERVER_PORT", "8080"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "myservice_db"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://gateway:8080/area_auth_api"),
		ServiceServiceURL: getEnv("SERVICE_SERVICE_URL", "http://gateway:8080/area_service_api"),
		AreaServiceURL:    getEnv("AREA_SERVICE_URL", "http://gateway:8080/area_area_api"),
		InternalSecret:    getEnv("INTERNAL_SECRET", ""),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
