package config

import (
	"os"
	"strconv"
)

type OAuth2Config struct {
	ClientID     string               `json:"client_id"`
	ClientSecret string               `json:"client_secret"`
	AuthURL      string               `json:"auth_url"`
	TokenURL     string               `json:"token_url"`
	RedirectURI  string               `json:"redirect_uri"`
	Scopes       []string             `json:"scopes"`
	UserInfoURL  string               `json:"user_info_url"`
	AuthParams   map[string]string    `json:"auth_params,omitempty"`
	Refresh      *OAuth2RefreshConfig `json:"refresh,omitempty"`
}

type OAuth2RefreshConfig struct {
	Enabled     bool              `json:"enabled"`
	TokenURL    string            `json:"token_url,omitempty"`
	Auth        string            `json:"auth,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

type FieldConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
}

type ProviderConfig struct {
	Name     string        `json:"name"`
	OAuth2   OAuth2Config  `json:"oauth2"`
	Mappings []FieldConfig `json:"mappings"`
}

type Config struct {
	HTTPPort                     string
	DBHost                       string
	DBPort                       string
	DBUser                       string
	DBPass                       string
	DBName                       string
	ServiceServiceURL            string
	InternalSecret               string
	OAuth2RefreshIntervalSeconds int
	OAuth2RefreshLeewayMinutes   int
}

func Load() Config {
	return Config{
		HTTPPort:                     getEnv("SERVER_PORT", "8080"),
		DBHost:                       getEnv("DB_HOST", "localhost"),
		DBPort:                       getEnv("DB_PORT", "5432"),
		DBUser:                       getEnv("DB_USER", "postgres"),
		DBPass:                       getEnv("DB_PASSWORD", "postgres"),
		DBName:                       getEnv("DB_NAME", "myservice_db"),
		ServiceServiceURL:            getEnv("SERVICE_SERVICE_URL", "http://gateway:8080/area_service_api"),
		InternalSecret:               getEnv("INTERNAL_SECRET", ""),
		OAuth2RefreshIntervalSeconds: getEnvInt("OAUTH2_REFRESH_INTERVAL_SECONDS", 60),
		OAuth2RefreshLeewayMinutes:   getEnvInt("OAUTH2_REFRESH_LEEWAY_MINUTES", 5),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}
	return def
}
