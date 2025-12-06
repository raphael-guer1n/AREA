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
}

type ProviderConfig struct {
	Name     string        `json:"name"`
	OAuth2   OAuth2Config  `json:"oauth2"`
	Mappings []FieldConfig `json:"mappings"`
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
