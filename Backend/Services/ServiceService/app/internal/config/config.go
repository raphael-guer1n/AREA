package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	HTTPPort string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPass   string
	DBName   string
}

func Load() Config {
	return Config{
		HTTPPort: getEnv("SERVER_PORT", "8080"),
		DBHost:   getEnv("DB_HOST", "localhost"),
		DBPort:   getEnv("DB_PORT", "5432"),
		DBUser:   getEnv("DB_USER", "postgres"),
		DBPass:   getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "myservice_db"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

type OAuth2Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	UserInfoURL  string   `json:"user_info_url"`
}

type MappingConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
}

type ProviderConfig struct {
	Name     string          `json:"name"`
	OAuth2   OAuth2Config    `json:"oauth2"`
	Mappings []MappingConfig `json:"mappings"`
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
	Mappings      []MappingConfig               `json:"mappings,omitempty"`
	Setup         *WebhookProviderSetupConfig `json:"setup,omitempty"`
	Teardown      *WebhookProviderSetupConfig `json:"teardown,omitempty"`
}

// LoadProviderConfigs loads all *.json provider configs from the given directory.
// Returns a map keyed by the provider name (e.g. "google", "discord").
func LoadProviderConfigs(dir string) (map[string]ProviderConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read providers dir: %w", err)
	}

	providers := make(map[string]ProviderConfig)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read provider file %s: %w", path, err)
		}

		var cfg ProviderConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal provider file %s: %w", path, err)
		}

		if cfg.Name == "" {
			return nil, fmt.Errorf("provider file %s: missing name", path)
		}

		for _, f := range cfg.Mappings {
			if f.FieldKey == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping missing field_key", cfg.Name)
			}
			if f.JSONPath == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping %s missing json_path", cfg.Name, f.FieldKey)
			}
			switch f.Type {
			case "string", "number", "boolean", "json":
			default:
				return nil, fmt.Errorf("provider %s: field %s has invalid type %q", cfg.Name, f.FieldKey, f.Type)
			}
		}

		providers[cfg.Name] = cfg
	}

	return providers, nil
}

type FieldConfig struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Label         string `json:"label"`
	Required      bool   `json:"required"`
	DefaultValuer string `json:"default"`
}

type OutputField struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

type ActionConfig struct {
	Title        string        `json:"title"`
	Label        string        `json:"label"`
	Type         string        `json:"type"`
	Fields       []FieldConfig `json:"fields"`
	OutputFields []OutputField `json:"output_fields"`
}

type ReactionConfig struct {
	Title      string          `json:"title"`
	Label      string          `json:"label"`
	Url        string          `json:"url"`
	Fields     []FieldConfig   `json:"fields"`
	Method     string          `json:"method"`
	BodyType   string          `json:"bodyType"`
	BodyStruct json.RawMessage `json:"body_struct"`
}

type ServiceConfig struct {
	Provider  string           `json:"provider"`
	Name      string           `json:"name"`
	IconURL   string           `json:"icon_url"`
	Label     string           `json:"label"`
	Actions   []ActionConfig   `json:"actions"`
	Reactions []ReactionConfig `json:"reactions"`
}

func LoadServiceConfig(dir string) (map[string]ServiceConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read services dir: %w", err)
	}

	services := make(map[string]ServiceConfig)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read service file %s: %w", path, err)
		}
		var cfg ServiceConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal service file %s: %w", path, err)
		}
		services[cfg.Name] = cfg
	}
	return services, nil
}

// LoadWebhookProviderConfigs loads all *.json webhook provider configs from the given directory.
// Returns a map keyed by the provider name (e.g. "github", "generic").
func LoadWebhookProviderConfigs(dir string) (map[string]WebhookProviderConfig, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read webhook providers dir: %w", err)
	}

	providers := make(map[string]WebhookProviderConfig)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read webhook provider file %s: %w", path, err)
		}

		var cfg WebhookProviderConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("unmarshal webhook provider file %s: %w", path, err)
		}

		if cfg.Name == "" {
			return nil, fmt.Errorf("webhook provider file %s: missing name", path)
		}

		if err := validateWebhookSignatureConfig(cfg.Name, cfg.Signature); err != nil {
			return nil, err
		}

		if err := validateWebhookProviderAction(cfg.Name, "setup", cfg.Setup); err != nil {
			return nil, err
		}
		if err := validateWebhookProviderAction(cfg.Name, "teardown", cfg.Teardown); err != nil {
			return nil, err
		}

		for _, f := range cfg.Mappings {
			if f.FieldKey == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping missing field_key", cfg.Name)
			}
			if f.JSONPath == "" {
				return nil, fmt.Errorf("webhook provider %s: mapping %s missing json_path", cfg.Name, f.FieldKey)
			}
			switch f.Type {
			case "string", "number", "boolean", "json":
			default:
				return nil, fmt.Errorf("webhook provider %s: field %s has invalid type %q", cfg.Name, f.FieldKey, f.Type)
			}
		}

		providers[cfg.Name] = cfg
	}

	return providers, nil
}

func validateWebhookProviderAction(providerName, label string, action *WebhookProviderSetupConfig) error {
	if action == nil {
		return nil
	}
	if action.Method == "" {
		return fmt.Errorf("webhook provider %s: %s method is required", providerName, label)
	}
	if action.URLTemplate == "" {
		return fmt.Errorf("webhook provider %s: %s url_template is required", providerName, label)
	}
	if action.Auth != nil {
		if action.Auth.Type == "" {
			return fmt.Errorf("webhook provider %s: %s auth type is required", providerName, label)
		}
		if action.Auth.Header == "" {
			return fmt.Errorf("webhook provider %s: %s auth header is required", providerName, label)
		}
		switch action.Auth.Type {
		case "oauth2":
		default:
			return fmt.Errorf("webhook provider %s: unsupported auth type %q", providerName, action.Auth.Type)
		}
	}
	if len(action.BodyTemplate) > 0 {
		var payload any
		if err := json.Unmarshal(action.BodyTemplate, &payload); err != nil {
			return fmt.Errorf("webhook provider %s: invalid %s body_template: %w", providerName, label, err)
		}
	}
	return nil
}

func validateWebhookSignatureConfig(providerName string, sig *WebhookSignatureConfig) error {
	if sig == nil {
		return nil
	}

	if sig.Type == "" {
		return fmt.Errorf("webhook provider %s: signature type is required", providerName)
	}

	signatureType, algorithm := normalizeSignatureType(sig)
	if signatureType == "token" {
		signatureType = "header"
	}

	switch signatureType {
	case "hmac":
		if sig.Header == "" {
			return fmt.Errorf("webhook provider %s: signature header is required", providerName)
		}
		if sig.SecretJSONPath == "" {
			return fmt.Errorf("webhook provider %s: signature secret_json_path is required", providerName)
		}
		if algorithm == "" {
			algorithm = "sha256"
		}
		switch algorithm {
		case "sha1", "sha256", "sha512":
		default:
			return fmt.Errorf("webhook provider %s: unsupported signature algorithm %q", providerName, algorithm)
		}
		if sig.Encoding != "" {
			switch strings.ToLower(sig.Encoding) {
			case "hex", "base64":
			default:
				return fmt.Errorf("webhook provider %s: unsupported signature encoding %q", providerName, sig.Encoding)
			}
		}
		if sig.TimestampToleranceSeconds < 0 {
			return fmt.Errorf("webhook provider %s: timestamp_tolerance_seconds must be >= 0", providerName)
		}
	case "header":
		if sig.Header == "" {
			return fmt.Errorf("webhook provider %s: signature header is required", providerName)
		}
		if sig.SecretJSONPath == "" {
			return fmt.Errorf("webhook provider %s: signature secret_json_path is required", providerName)
		}
	default:
		return fmt.Errorf("webhook provider %s: unsupported signature type %q", providerName, sig.Type)
	}

	return nil
}

func normalizeSignatureType(sig *WebhookSignatureConfig) (string, string) {
	signatureType := strings.ToLower(sig.Type)
	algorithm := strings.ToLower(sig.Algorithm)

	if strings.HasPrefix(signatureType, "hmac-") {
		algorithm = strings.TrimPrefix(signatureType, "hmac-")
		signatureType = "hmac"
	}

	return signatureType, algorithm
}
