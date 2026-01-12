package oauth2

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// ProviderConfig defines configuration for a single OAuth2 provider
type ProviderConfig struct {
	Name         string            `json:"name"`
	ClientID     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	AuthURL      string            `json:"auth_url"`
	TokenURL     string            `json:"token_url"`
	RedirectURI  string            `json:"redirect_uri"`
	Scopes       []string          `json:"scopes"`
	UserInfoURL  string            `json:"user_info_url"`
	AuthParams   map[string]string `json:"auth_params,omitempty"`
	Refresh      *RefreshConfig    `json:"refresh,omitempty"`
}

type RefreshConfig struct {
	Enabled     bool              `json:"enabled"`
	TokenURL    string            `json:"token_url,omitempty"`
	Auth        string            `json:"auth,omitempty"`
	ContentType string            `json:"content_type,omitempty"`
	Params      map[string]string `json:"params,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// ConfigLoader manages to load OAuth2 configurations from service-service API
type ConfigLoader struct {
	serviceServiceURL string
	internalSecret    string
	httpClient        *http.Client
	cache             map[string]*ProviderConfig
	cacheMutex        sync.RWMutex
	serviceList       []string
	serviceListTime   time.Time
}

// NewConfigLoader creates a new config loader that fetches from service-service API
func NewConfigLoader(serviceServiceURL string, internalSecret string) *ConfigLoader {
	return &ConfigLoader{
		serviceServiceURL: serviceServiceURL,
		internalSecret:    internalSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: make(map[string]*ProviderConfig),
	}
}

// GetProvider fetches a provider configuration from service-service API (with caching)
func (l *ConfigLoader) GetProvider(name string) (*ProviderConfig, error) {
	// Check cache first
	l.cacheMutex.RLock()
	if cached, exists := l.cache[name]; exists {
		l.cacheMutex.RUnlock()
		return cached, nil
	}
	l.cacheMutex.RUnlock()

	// Fetch from API
	l.cacheMutex.Lock()
	defer l.cacheMutex.Unlock()

	// Double-check after acquiring write lock
	if cached, exists := l.cache[name]; exists {
		return cached, nil
	}

	// Fetch OAuth2 config from service-service
	url := fmt.Sprintf("%s/providers/oauth2-config?service=%s", l.serviceServiceURL, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider config request: %w", err)
	}
	if strings.TrimSpace(l.internalSecret) != "" {
		req.Header.Set("X-Internal-Secret", l.internalSecret)
	}
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch provider config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service-service returned status %d for provider %s", resp.StatusCode, name)
	}

	var apiResp struct {
		Success bool `json:"success"`
		Data    struct {
			ClientID     string            `json:"client_id"`
			ClientSecret string            `json:"client_secret"`
			AuthURL      string            `json:"auth_url"`
			TokenURL     string            `json:"token_url"`
			RedirectURI  string            `json:"redirect_uri"`
			Scopes       []string          `json:"scopes"`
			UserInfoURL  string            `json:"user_info_url"`
			AuthParams   map[string]string `json:"auth_params,omitempty"`
			Refresh      *RefreshConfig    `json:"refresh,omitempty"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("service-service error: %s", apiResp.Error)
	}

	// Convert to ProviderConfig and resolve environment variables
	config := &ProviderConfig{
		Name:         name,
		ClientID:     resolveEnvVar(apiResp.Data.ClientID),
		ClientSecret: resolveEnvVar(apiResp.Data.ClientSecret),
		AuthURL:      apiResp.Data.AuthURL,
		TokenURL:     apiResp.Data.TokenURL,
		RedirectURI:  apiResp.Data.RedirectURI,
		Scopes:       apiResp.Data.Scopes,
		UserInfoURL:  apiResp.Data.UserInfoURL,
		AuthParams:   apiResp.Data.AuthParams,
		Refresh:      apiResp.Data.Refresh,
	}

	// Cache it
	l.cache[name] = config

	return config, nil
}

// resolveEnvVar resolves an environment variable name to its value
// If the value looks like an env var name (uppercase with underscores), it tries to resolve it
// Otherwise returns the value as-is
func resolveEnvVar(value string) string {
	// If value is empty, return as-is
	if value == "" {
		return value
	}

	// Try to get from the environment
	if envValue := os.Getenv(value); envValue != "" {
		return envValue
	}

	// If not found in env, return the original value (could be a literal value)
	return value
}

// ListProviders fetches the list of available providers from service-service API
func (l *ConfigLoader) ListProviders() ([]string, error) {
	// Cache service list for 1 minute
	l.cacheMutex.RLock()
	if time.Since(l.serviceListTime) < time.Minute && l.serviceList != nil {
		list := l.serviceList
		l.cacheMutex.RUnlock()
		return list, nil
	}
	l.cacheMutex.RUnlock()

	// Fetch from API
	url := fmt.Sprintf("%s/providers/services", l.serviceServiceURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create service list request: %w", err)
	}
	if strings.TrimSpace(l.internalSecret) != "" {
		req.Header.Set("X-Internal-Secret", l.internalSecret)
	}
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service list: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("service-service returned status %d", resp.StatusCode)
	}

	var apiResp struct {
		Success bool `json:"success"`
		Data    struct {
			Services []string `json:"services"`
		} `json:"data"`
		Error string `json:"error,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("service-service error: %s", apiResp.Error)
	}

	// Cache the list
	l.cacheMutex.Lock()
	l.serviceList = apiResp.Data.Services
	l.serviceListTime = time.Now()
	l.cacheMutex.Unlock()

	return apiResp.Data.Services, nil
}
