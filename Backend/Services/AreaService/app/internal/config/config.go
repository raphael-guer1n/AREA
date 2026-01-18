package config

import (
	"os"
	"strings"
)

type Config struct {
	HTTPPort              string
	DBHost                string
	DBPort                string
	DBUser                string
	DBPass                string
	DBName                string
	AuthServiceURL        string
	ServiceServiceURL     string
	AreaServiceURL        string
	InternalSecret        string
	CreateActionsUrls     map[string]string
	DelActionsUrls        map[string]string
	ActivateActionsUrls   map[string]string
	DeactivateActionsUrls map[string]string
}

func Load() Config {
	createActionsUrls := GetActionsUrls("CREATE_ACTIONS_URLS")
	delActionsUrls := GetActionsUrls("DEL_ACTIONS_URLS")
	activateActionsUrls := GetActionsUrls("ACTIVATE_ACTIONS_URLS")
	deactivateActionsUrls := GetActionsUrls("DEACTIVATE_ACTIONS_URLS")

	return Config{
		HTTPPort:              getEnv("SERVER_PORT", "8080"),
		DBHost:                getEnv("DB_HOST", "localhost"),
		DBPort:                getEnv("DB_PORT", "5432"),
		DBUser:                getEnv("DB_USER", "postgres"),
		DBPass:                getEnv("DB_PASSWORD", "postgres"),
		DBName:                getEnv("DB_NAME", "myservice_db"),
		AuthServiceURL:        getEnv("AUTH_SERVICE_URL", "http://gateway:8080/area_auth_api"),
		ServiceServiceURL:     getEnv("SERVICE_SERVICE_URL", "http://gateway:8080/area_service_api"),
		AreaServiceURL:        getEnv("AREA_SERVICE_URL", "http://gateway:8080/area_area_api"),
		InternalSecret:        getEnv("INTERNAL_SECRET", ""),
		CreateActionsUrls:     createActionsUrls,
		DelActionsUrls:        delActionsUrls,
		ActivateActionsUrls:   activateActionsUrls,
		DeactivateActionsUrls: deactivateActionsUrls,
	}
}

func GetActionsUrls(envVarName string) map[string]string {
	urls := make(map[string]string)
	data := os.Getenv(envVarName)

	if data == "" {
		return urls
	}

	data = strings.TrimSpace(data)
	data = strings.Trim(data, "{}")

	pairs := strings.Split(data, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		colonIdx := strings.Index(pair, ":")
		if colonIdx == -1 {
			continue
		}

		key := strings.Trim(strings.TrimSpace(pair[:colonIdx]), `"`)
		value := strings.Trim(strings.TrimSpace(pair[colonIdx+1:]), `"`)

		if key != "" {
			urls[key] = value
		}
	}
	return urls
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
