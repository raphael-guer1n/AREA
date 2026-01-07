package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type FieldConfig struct {
	FieldKey string `json:"field_key"`
	JSONPath string `json:"json_path"`
	Type     string `json:"type"`
	Optional bool   `json:"optional,omitempty"`
}

type ProviderConfig struct {
	Name     string        `json:"name"`
	OAuth2   OAuth2Config  `json:"oauth2"`
	Mappings []FieldConfig `json:"mappings"`
}

type WebhookSignatureConfig struct {
	Type           string `json:"type"`
	Header         string `json:"header"`
	Prefix         string `json:"prefix"`
	SecretJSONPath string `json:"secret_json_path"`
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

		if cfg.Signature != nil {
			if cfg.Signature.Type == "" {
				return nil, fmt.Errorf("webhook provider %s: signature type is required", cfg.Name)
			}
			if cfg.Signature.Header == "" {
				return nil, fmt.Errorf("webhook provider %s: signature header is required", cfg.Name)
			}
			if cfg.Signature.SecretJSONPath == "" {
				return nil, fmt.Errorf("webhook provider %s: signature secret_json_path is required", cfg.Name)
			}

			switch cfg.Signature.Type {
			case "hmac-sha256", "hmac-sha1":
			default:
				return nil, fmt.Errorf("webhook provider %s: unsupported signature type %q", cfg.Name, cfg.Signature.Type)
			}
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
