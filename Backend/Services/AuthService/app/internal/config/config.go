package config

import "os"

type Config struct {
	HTTPPort         string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPass           string
	DBName           string
	OAuth2ConfigPath string
}

func Load() Config {
	return Config{
		HTTPPort:         getEnv("SERVER_PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPass:           getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "myservice_db"),
		OAuth2ConfigPath: getEnv("OAUTH2_CONFIG_PATH", "oauth2.config.json"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
