package oauth2

import (
	"encoding/json"
	"fmt"
	"os"
)

// ProviderConfig defines configuration for a single OAuth2 provider
type ProviderConfig struct {
	Name         string   `json:"name"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	RedirectURI  string   `json:"redirect_uri"`
	Scopes       []string `json:"scopes"`
	UserInfoURL  string   `json:"user_info_url"`
}

// ConfigOAuth2 holds all OAuth2 provider configurations
type ConfigOAuth2 struct {
	Providers map[string]ProviderConfig `json:"providers"`
}

// LoadConfig loads OAuth2 configuration from a JSON file
func LoadConfig(path string) (*ConfigOAuth2, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config ConfigOAuth2
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	// Validate that we have at least one provider
	if len(config.Providers) == 0 {
		return nil, fmt.Errorf("no OAuth2 providers configured")
	}

	return &config, nil
}

// GetProvider returns a specific provider configuration
func (c *ConfigOAuth2) GetProvider(name string) (*ProviderConfig, error) {
	provider, ok := c.Providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return &provider, nil
}

// ListProviders returns all available provider names
func (c *ConfigOAuth2) ListProviders() []string {
	providers := make([]string, 0, len(c.Providers))
	for name := range c.Providers {
		providers = append(providers, name)
	}
	return providers
}
