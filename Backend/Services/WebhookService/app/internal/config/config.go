package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	HTTPPort          string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPass            string
	DBName            string
	ServiceServiceURL string
	AuthServiceURL    string
	AreaServiceURL    string
	PublicBaseURL     string
}

func Load() Config {
	return Config{
		HTTPPort:          getEnv("SERVER_PORT", "8085"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPass:            getEnv("DB_PASSWORD", "postgres"),
		DBName:            getEnv("DB_NAME", "webhook_service_db"),
		ServiceServiceURL: getEnv("SERVICE_SERVICE_URL", "http://gateway:8080/service-service"),
		AuthServiceURL:    getEnv("AUTH_SERVICE_URL", "http://gateway:8080/auth-service"),
		AreaServiceURL:    getEnv("AREA_SERVICE_URL", "http://gateway:8080/area-service"),
		PublicBaseURL:     getEnv("PUBLIC_BASE_URL", "http://gateway:8080/webhook-service"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type FieldConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
}

type WebhookSignatureConfig struct {
	Type                      string `json:"type"`
	Header                    string `json:"header"`
	Prefix                    string `json:"prefix"`
	SecretJSONPath            string `json:"secret_json_path"`
	Algorithm                 string `json:"algorithm,omitempty"`
	Encoding                  string `json:"encoding,omitempty"`
	SigningStringTemplate     string `json:"signing_string_template,omitempty"`
	TimestampHeader           string `json:"timestamp_header,omitempty"`
	TimestampToleranceSeconds int    `json:"timestamp_tolerance_seconds,omitempty"`
}

type WebhookProviderAuthConfig struct {
	Type   string `json:"type"`
	Header string `json:"header"`
	Prefix string `json:"prefix"`
}

type WebhookProviderSetupConfig struct {
	Method             string                     `json:"method"`
	URLTemplate        string                     `json:"url_template"`
	Headers            map[string]string          `json:"headers,omitempty"`
	Auth               *WebhookProviderAuthConfig `json:"auth,omitempty"`
	BodyTemplate       json.RawMessage            `json:"body_template,omitempty"`
	ResponseIDJSONPath string                     `json:"response_id_json_path"`
}

type WebhookProviderConfig struct {
	Name          string                      `json:"name"`
	Signature     *WebhookSignatureConfig     `json:"signature,omitempty"`
	EventHeader   string                      `json:"event_header"`
	EventJSONPath string                      `json:"event_json_path"`
	Mappings      []FieldConfig               `json:"mappings,omitempty"`
	Setup         *WebhookProviderSetupConfig `json:"setup,omitempty"`
	Teardown      *WebhookProviderSetupConfig `json:"teardown,omitempty"`
}
